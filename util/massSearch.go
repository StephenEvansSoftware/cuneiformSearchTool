package util

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
)

const (
	massSearchFileNamePrefix = "cst-"
	massSearchFileType       = ".csv"

	// TODO: make this configurable
	massSearchOutputFolder = "massSearchOutputs"
)

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func GetTotalMassSearches(termsLength int, properNounGrouping bool) int {
	baseN := totalTermsWithoutProperNounGrouping
	if properNounGrouping {
		baseN = totalTermsWithProperNounGrouping
	}
	return powInt(baseN, termsLength)
}

/*
Gathers settings and creates a hash
Creates the file with that hash if it doesn't already exist
Loads the file and checks the settings match
Ensures the corpus selection matches
Creates all possible permutations of valid types
Iterates through the file checking which permutations have been finished
Continues searching with the remaining permutations
*/
func RunMassSearch(searchSettings *SearchSettings, massSearchSettings *MassSearchSettings) {
	// Gather settings and get first 8 characters of the hash of the set
	settings := getMassSearchSettingStrings(searchSettings, massSearchSettings)
	sortedSettings := strings.Join(settings, "\n")
	settingsHash := hash([]byte(sortedSettings))[:8]

	// Check if file exists already
	fileName := fmt.Sprintf("%s%s%s", massSearchFileNamePrefix, settingsHash, massSearchFileType)
	filePath := fmt.Sprintf("%s/%s", massSearchOutputFolder, fileName)
	isDir, isFile, err := checkExists(filePath)
	if err != nil {
		panic(err) // TODO: handle better
	}
	if isDir {
		panic(fmt.Errorf("%s is a directory, not a file", filePath)) // TODO: handle better
	}
	if !isFile {
		// No file found, create it now
		corpusLine := fmt.Sprintf(",%s", strings.Join(getCorpusList(searchSettings.SearchAllTexts, searchSettings.SelectedCorpus), ","))
		toWrite := fmt.Sprintf("%s\n\n%s", sortedSettings, corpusLine)
		err = os.WriteFile(filePath, []byte(toWrite), 0755)
		if err != nil {
			panic(err) // TODO: handle better
		}
	}

	// Load file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0666) // TODO TODO TODO
	if err != nil {
		panic(err) // TODO: handle better
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	// Ensure settings listed in the file are the same
	for _, s := range settings {
		fileScanner.Scan()
		currLine := fileScanner.Text()
		if currLine != s {
			panic(fmt.Errorf("'%s' != '%s'", currLine, s)) // TODO: handle better
		}
	}
	// Settings are all correct. Ensure no additional settings
	fileScanner.Scan()
	if fileScanner.Text() != "" {
		panic(fmt.Errorf("expecting newline after settings in file")) // TODO: handle better
	}
	// Last check to make sure that the first line is the corpus list
	// TODO: probably don't need to do this at this point. once settings are fine assume data is too
	fileScanner.Scan()
	if corpusLine := fileScanner.Text(); corpusLine != fmt.Sprintf(",%s", strings.Join(getCorpusList(searchSettings.SearchAllTexts, searchSettings.SelectedCorpus), ",")) {
		panic(fmt.Errorf("unknown corpus line: %s", corpusLine))
	}

	// Create the list of all possible term permutations
	searchTerms := ValidTypes
	if massSearchSettings.UseProperNounGrouping {
		searchTerms = append([]string{AllProperNounsType}, validTypesSansProperNouns...)
	}
	possibleTerms := len(searchTerms)
	permutationIndexes := [][]int{}

	for i := 0; i < GetTotalMassSearches(massSearchSettings.TermsLength, massSearchSettings.UseProperNounGrouping); i++ {
		permutation := []int{}
		for j := massSearchSettings.TermsLength - 1; j >= 0; j-- {
			e := i / powInt(possibleTerms, j) % possibleTerms
			permutation = append(permutation, e)
		}
		permutationIndexes = append(permutationIndexes, permutation)
	}

	// Iterate through the file and check the first entry is the same as the current constructed term
	completedSearches := 0
	for fileScanner.Scan() {
		completedSearches++
		line := fileScanner.Text()
		searched := strings.Split(line, ",")[0]
		currSearch := [][]string{}
		for _, i := range permutationIndexes[completedSearches-1] {
			currSearch = append(currSearch, []string{searchTerms[i]})
		}
		//		currSearch := fmt.Sprintf("%v", permutationIndexes[completedSearches - 1])
		if searched != fmt.Sprintf("%v", currSearch) {
			panic(fmt.Errorf("'%s' search does not match '%v'", searched, currSearch))
			// TODO: handle better
		}
	}

	// Perform our search here
	for i := completedSearches; i < len(permutationIndexes); i++ {
		terms := [][]string{}
		for j := 0; j < massSearchSettings.TermsLength; j++ {
			terms = append(terms, []string{searchTerms[permutationIndexes[i][j]]})
		}
		newSearchResultLine := fmt.Sprintf("\n%v", terms)
		// Actual search happening here
		csr, _, err := SearchFor(terms, searchSettings)
		if err != nil {
			panic(err.Error()) // TODO: handle better
		}
		// Create the new line in the search results file
		for _, c := range csr {
			totalCorpusResults := 0
			for _, t := range c.TextResults {
				totalCorpusResults += len(t.Results)
			}
			newSearchResultLine = fmt.Sprintf("%s,%d", newSearchResultLine, totalCorpusResults)
		}
		// Write to open file
		_, err = file.Write([]byte(newSearchResultLine))
		if err != nil {
			panic(err.Error()) // TODO: handle better
		}
	}
}

func getMassSearchSettingStrings(ss *SearchSettings, mss *MassSearchSettings) []string {
	// Save settings as string that we can sort and hash later
	settings := []string{fmt.Sprintf("TermsLength: %d", mss.TermsLength)}
	if ss.SearchAllTexts {
		settings = append(settings, "ScopeAllCorpus")
	} else {
		selectedCorpus := []string{}
		for _, c := range ss.SelectedCorpus {
			selectedCorpus = append(selectedCorpus, c)
		}
		slices.Sort(selectedCorpus)
		settings = append(settings, fmt.Sprintf("ScopeSelectedCorpus: %v", selectedCorpus))
	}
	if ss.SearchWrappedText {
		settings = append(settings, "WrapLines")
	}
	if ss.UsePOS {
		settings = append(settings, "UsePOS")
	}
	if mss.UseProperNounGrouping {
		settings = append(settings, "UseProperNounGrouping")
	}
	slices.Sort(settings)
	return settings
}

func hash(toHash []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(toHash))
}
