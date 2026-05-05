package util

import (
	"bufio"
	"crypto/sha256"
	"fmt"
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

func GetTotalMassSearches(termsLength int, properNounGrouping bool) int {
	baseN := totalTermsWithoutProperNounGrouping
	if properNounGrouping {
		baseN = totalTermsWithProperNounGrouping
	}
	// Avoiding importing the maths package to do a simple pow
	output := baseN
	for i := 1; i < termsLength; i++ {
		output *= baseN
	}
	return output
}

func RunMassSearch(searchSettings *SearchSettings, massSearchSettings *MassSearchSettings) {
	/*
				TODO: a lot...
		1. using the settings, create a hash
		2. check if a file exists with that hash included in its name
		3. load file
		4. check if file contains the settings as they are in the current search
		4a. yes: we have found our previous attempt at searching
							count the number of entries
		4b. no: we create the file with our hash and search and mass search settings
				5. compute all possible permutations of the types depending on settings
					note: this will be n^r
						n = total types being used (all Valid types or validSansProperNouns + 1)
						r = term length
					5a. skip the first X terms where X is the number of current entires (if applicable)
				6. start the search
					6a. each successful search we append to file
						ordering of corpus results is alphabitical
						format:
									c1	c2	c3
							[a,a,a]	x	y	z
							[a,a,b]	x2	y2	z2
						will use csv to represent, so ends up looking like

							,c1,c2,c3
							"[a,a,a]",x,y,z
							"[a,a,b]",x2,y2,z2
	*/

	/*
		1. create hash of settings
			settings will be simple strings of things that ARE set
	*/
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
	file, err := os.Open(filePath)
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

	/*
		file format:

		Ordered
		ListOf
		Settings
					< newline break between settings and data
		,c1,c2,c3
		[a a a],12,99,41
		[a a b],82,14,55
	*/

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
