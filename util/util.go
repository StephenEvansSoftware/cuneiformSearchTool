package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"cuneiformSearchTool/text"
)

var (
	// Valid word types
	ValidTypes = []string{
		"AJ",	// Adjective
		"AV",	// Adverb
		"CN",	// Celestial name
		"CNJ",	// Conjunction
		"DET",	// Determinative pronoun
		"DN",	// Divine name
		"DP",	// Demonstrative pronoun
		"EN",	// Ethnic name
		"GN",	// Geographic name
		"IP",	// Independant pronoun
		"J",	// Exclamation / interjection
		"MN",	// Month name
		"MOD",	// Modifier
		"N",	// Noun
		"NU",	// Numeral
		"ON",	// Settlement name
		"PN",	// Personal name
		"PP",	// Possesive pronoun
		"PRP",	// Preposition
		"QN",	// Place name
		"QP",	// Question particle
		"REL",	// Relative pronoun
		"RN",	// Ruler name
		"RP",	// Reflexive prnoun
		"SBJ",	// Subjunction
		"SN",	// City / settlement name
		"TN",	// Temple name
		"V",	// Verb
		"V/i",	//// Some sumerian stuff
		"V/t",	//// Some sumerian stuff
		"WN",	// Water name
		"X",	// Meaning unknown
		"XP",	// Indefinite pronoun
		"n",	//// Breaks? Only on POS
		"u",	//// Breaks? Only on POS
	}

	properNounTypes = []string{
		"CN",
		"DN",
		"EN",
		"GN",
		"MN",
		"ON",
		"PN",
		"QN",
		"RN",
		"SN",
		"TN",
		"WN",
	}

	inputSources = "./jsonFiles/sources"
	inputCorpus  = "./jsonFiles/corpus"

	corpusTextInfo = []*CorpusInfo{}
)

type TextSearchResult struct {
	TextName string
	Results  []SearchResult
	TextPath string // Used to pass back to load if selected
}

type SearchResult struct {
	OverallLineIndex int
	Label            string
	Position         int
	LastEntry        string
}

type CorpusInfo struct {
	CorpusName string
	CorpusPath string
	Texts      []*TextInfo
}

type TextInfo struct {
	TextName string
	TextPath string
}

// TODO TODO TODO: change this a LOT
type SearchSettings struct {
	SearchAllTexts bool
	SelectedCorpus []string

	SearchWrappedText        bool
	PrintTextMetaInformation bool
	UsePOS                   bool
	UseFORM                  bool

	// TODO: legacy, not used
	printWhileSearching bool
}

// Package intiialisation
func init() {
	if err := loadCorpusTexts(); err != nil {
		panic(err.Error())
	}
}

// Load corpus names and text paths
func loadCorpusTexts() error {
	// Check corpus dir exists
	isDir, _, err := checkExists(inputCorpus)
	if err != nil {
		return err
	}
	if !isDir {
		return fmt.Errorf("corpus directory '%s' not found", inputCorpus)
	}
	// Load the corpus directory
	corpusDirs, err := os.ReadDir(inputCorpus)
	if err != nil {
		return err
	}
	// Iterate through the entries, creating new corpus for each dir
	for _, corpus := range corpusDirs {
		corpusPath := fmt.Sprintf("%s/%s", inputCorpus, corpus.Name())
		// Ensure the corpus is a directory
		corpusStat, _ := os.Stat(corpusPath)
		if !corpusStat.IsDir() {
			continue
			// TODO: handle floating files better, for now just move on
		}
		corpusInfo := &CorpusInfo{
			CorpusName: corpus.Name(),
			CorpusPath: corpusPath,
		}
		texts := []*TextInfo{}
		// Iterate through the texts contained in the corpus
		corpusTexts, _ := os.ReadDir(corpusPath)
		for _, text := range corpusTexts {
			textPath := fmt.Sprintf("%s/%s", corpusPath, text.Name())
			// Ensure text is a regular file
			textStat, _ := os.Stat(textPath)
			switch textMode := textStat.Mode(); {
			case textMode.IsDir():
				continue // TODO: handle floating dirs better, for now just move on
			case textMode.IsRegular():
				// Add the text to the list under the current corpus
				texts = append(texts, &TextInfo{
					TextName: text.Name(),
					TextPath: textPath,
				})
			default:
				continue // TODO: unsure if needed, just leaving here as a placeholder for now
			}
		}
		corpusInfo.Texts = texts
		corpusTextInfo = append(corpusTextInfo, corpusInfo)
	}
	return nil
}

func ListCorpusAndTexts() []*CorpusInfo {
	return corpusTextInfo
}

func ValidateTypes(types [][]string) error {
	invalidTypes := []string{}
	if len(types) == 0 {
		return fmt.Errorf("no types found")
	}
	for _, searchArr := range types {
		if len(searchArr) == 0 {
			invalidTypes = append(invalidTypes, "empty term invalid")
		}
		for _, search := range searchArr {
			typeFound := false
			for _, t := range ValidTypes {
				if search == t {
					typeFound = true
					break
				}
			}
			if !typeFound {
				invalidTypes = append(invalidTypes, search)
			}
		}
	}
	if len(invalidTypes) != 0 {
		return fmt.Errorf("invalid types found: %+v", invalidTypes)
	}
	return nil
}

// TODO: change inputs and a whole lot here
// TODO TODO TODO: the gui should send the information back here of which ones to search and their paths and such
// TODO TODO TODO TODO TODO: just accept search settings instead of all the extra things
func SearchFor(terms [][]string, searchSettings *SearchSettings) ([]*TextSearchResult, int, error) {
	searchResults := []*TextSearchResult{}
	totalInstancesFound := 0
	// When allCorpus is set, change which corpus to search by overriding its value
	if searchSettings.SearchAllTexts {
		searchSettings.SelectedCorpus = []string{}
		for _, corp := range corpusTextInfo {
			searchSettings.SelectedCorpus = append(searchSettings.SelectedCorpus, corp.CorpusName)
		}
	}

	for _, corpTextInfo := range corpusTextInfo {
		for _, corp := range searchSettings.SelectedCorpus {
			if corp == corpTextInfo.CorpusName {
				// Search here
				for _, textInfo := range corpTextInfo.Texts {
					// TODO TODO TODO: add text deliniation here
					textFile, err := readFile(textInfo.TextPath)
					if err != nil {
						return nil, -1, err
					}
					searchResult, err := searchText(textFile, terms, searchSettings)

					if err != nil {
						return nil, -1, err
					}
					if len(searchResult.Results) != 0 {
						searchResult.TextPath = textInfo.TextPath
						searchResults = append(searchResults, searchResult)
						totalInstancesFound += len(searchResult.Results)
					}
				}
			}
		}
	}
	return searchResults, totalInstancesFound, nil
}

func checkExists(path string) (isDir, isFile bool, err error) {
	fi, err := os.Stat(path)
	if err == nil {
		mode := fi.Mode()
		return mode.IsDir(), mode.IsRegular(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, false, nil
	}
	return false, false, err
}

func readFile(path string) (*text.Text, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileBytes, _ := ioutil.ReadAll(file)

	var text text.Text
	if err := json.Unmarshal(fileBytes, &text); err != nil {
		return nil, err
	}
	return &text, nil
}

// TODO: a lot. this was just taken pretty much verbatim from my initial prototype god-file monstrosity
func searchText(textBody *text.Text, searchPattern [][]string, searchSettings *SearchSettings) (*TextSearchResult, error) {
	if len(searchPattern) == 0 {
		return nil, fmt.Errorf("cannot have empty search paattern")
	}

	// Simple print function to avoid a bunch of if statements in the code
	printFunc := func(s string, args ...interface{}) {} // Instantiate as doing nothing
	if searchSettings.printWhileSearching {
		printFunc = func(s string, args ...interface{}) {
			fmt.Printf(s, args...)
		}
	}

	spLen := len(searchPattern)
	previouslyFoundPatternPart := make([]bool, spLen)
	foundPatterns := &TextSearchResult{
		TextName: textBody.TextId,
		Results:  []SearchResult{},
	}
	fragIndex := 0
	currLabel := ""
	// TODO: move this somewhere else
	metaInf := &text.TextMetaInfo{
		FragCounts: map[string]int{},
	}

	printFunc("Start of text: %+v", textBody.TextId)
	numLines := 1

	for _, cdlOne := range textBody.Cdl {
		for _, cdlTwo := range cdlOne.Cdl {
			if cdlTwo.Subtype == "body" {
				// Here we have the actual body of the text
				for _, cdlThree := range cdlTwo.Cdl {
					if len(cdlThree.Cdl) != 0 {
						for _, cdlFour := range cdlThree.Cdl {
							var cdlChecker func(cdl4 *text.CdlLevelFour)
							cdlChecker = func(cdl4 *text.CdlLevelFour) {
								if cdl4.Type == "line-start" {
									//printFunc("\n\tLine %d\n\t\t", numLines)
									printFunc("\n\t%d: %s\n\t\t", numLines, cdl4.Label)
									currLabel = cdl4.Label
									numLines += 1
									fragIndex = 0
									// Check if we are searching across multiple lines
									if !searchSettings.SearchWrappedText {
										// Simply reset the previouslyFoundPatternPart by remaking it
										previouslyFoundPatternPart = make([]bool, spLen)
									}
								}
								if cdl4.F != nil {
									fragIndex += 1
									// This is the actual element in the text we care about
									printFunc("%s ", cdl4.Frag)
									metaInf.TotalFrags += 1
									//metaInf.FragCounts[cdl4.F.Epos] = metaInf.FragCounts[cdl4.F.Epos] + 1
									if searchSettings.UsePOS {
										metaInf.FragCounts[cdl4.F.Pos] += 1
									} else {
										metaInf.FragCounts[cdl4.F.Epos] += 1
									}
									// Reverse iterate through from element N-1 to 0
									for pfppi := spLen - 2; pfppi >= 0; pfppi -= 1 {
										// Checking position n
										// If position n set, we have found pattern up to n
										// Check n+1 and set accordingly
										// Then unset position n so n-1 can continue and set accordingly
										if previouslyFoundPatternPart[pfppi] {
											// Loop through n+1 searchPattern to see if any of them match the epos (or pos)
											for _, spf := range searchPattern[pfppi+1] {
												if searchSettings.UsePOS {
													if cdl4.F.Pos == spf {
														// Pattern part found, set found on increment to know to check / accept pattern
														previouslyFoundPatternPart[pfppi+1] = true
														break
													}
												} else {
													if cdl4.F.Epos == spf {
														// Pattern part found, set found on increment to know to check / accept pattern
														previouslyFoundPatternPart[pfppi+1] = true
														break
													}
												}
											}
											// Unset pfpp[pfppi] as we've finished checking pfppi+1
											previouslyFoundPatternPart[pfppi] = false
										}
									}
									// Always have to check current fragment against first search pattern
									previouslyFoundPatternPart[0] = false // Safety check to make sure if we are searching for a single element we don't always have 0 set
									for _, spf := range searchPattern[0] {
										if searchSettings.UsePOS {
											if cdl4.F.Pos == spf {
												previouslyFoundPatternPart[0] = true
											}
											break
										} else {
											if cdl4.F.Epos == spf {
												previouslyFoundPatternPart[0] = true
											}
											break
										}
									}
									// Final check to see if we have completed the pattern or not
									if previouslyFoundPatternPart[spLen-1] {
										// WE HAVE FOUND THE PATTERN
										// Store line and number
										newResult := SearchResult{
											OverallLineIndex: numLines - 1,
											Label:            currLabel,
											Position:         fragIndex,
										}
										if searchSettings.UseFORM {
											newResult.LastEntry = cdl4.F.Form
										} else {
											newResult.LastEntry = cdl4.Frag
										}
										foundPatterns.Results = append(foundPatterns.Results, newResult)
										previouslyFoundPatternPart[spLen-1] = false
									}
								}
								if cdl4.Cdl != nil {
									for _, cdlFive := range cdl4.Cdl {
										cdlChecker(cdlFive)
									}
								}
							}
							cdlChecker(cdlFour)
						}
					}
				}
			}
		}
	}
	printFunc("\nEnd of text")

	metaInf.TotalLines = numLines
	metaInf.LineLengthAverage = float32(metaInf.TotalFrags) / float32(numLines)
	metaInf.TotalFragsWithoutBlanks = metaInf.TotalFrags - metaInf.FragCounts[""]
	metaInf.LineLengthAverageWithoutBlanks = float32(metaInf.TotalFragsWithoutBlanks) / float32(numLines)

	// TODO: during the search, we don't want to print anything. change / remove these
	//	if printTextMetaInfo {
	if false {
		fmt.Printf("Meta information on the text:\n"+
			"\tTotalFrags: %d\n"+
			"\tWithout Xs: %d\n"+
			"\tTotalLines: %d\n"+
			"\tLineLengthAverage: %.2f\n"+
			"\tAverageWithout Xs: %.2f\n",
			metaInf.TotalFrags,
			metaInf.TotalFragsWithoutBlanks,
			metaInf.TotalLines,
			metaInf.LineLengthAverage,
			metaInf.LineLengthAverageWithoutBlanks,
		)
		if metaInf.TotalFrags != 0 {
			fmt.Println("\tFragCounts:")
			for s, i := range metaInf.FragCounts {
				fmt.Printf("\t\t%s\t%d\n", s, i)
			}
		}
	}

	return foundPatterns, nil
}

func GetPrintableText(textPath string, metaInf bool, usePOS bool, useFORM bool) (string, error) {
	text, err := readFile(textPath)
	if err != nil {
		return "", err
	}
	return getPrintableText(text, metaInf, usePOS, useFORM), nil
}

func getPrintableText(textBody *text.Text, printMetaInfo bool, usePOS bool, useFORM bool) string {
	textReturn := fmt.Sprintf("Start of text: %+v\n", textBody.TextId)
	numLines := 1

	metaInf := &text.TextMetaInfo{
		FragCounts: map[string]int{},
	}

	for _, cdlOne := range textBody.Cdl {
		for _, cdlTwo := range cdlOne.Cdl {
			if cdlTwo.Subtype == "body" {
				// Here we have the actual body of the text
				for _, cdlThree := range cdlTwo.Cdl {
					if len(cdlThree.Cdl) != 0 {
						for _, cdlFour := range cdlThree.Cdl {
							var cdlChecker func(cdl4 *text.CdlLevelFour)
							cdlChecker = func(cdl4 *text.CdlLevelFour) {
								if cdl4.Type == "line-start" {
									textReturn = fmt.Sprintf("%s\n\t%d: %s\n\t\t", textReturn, numLines, cdl4.Label)
									numLines += 1
								}
								if cdl4.F != nil {
									// This is the actual element in the text we care about
									if useFORM {
										textReturn = fmt.Sprintf("%s%s ", textReturn, cdl4.F.Form)
									} else {
										textReturn = fmt.Sprintf("%s%s ", textReturn, cdl4.Frag)
									}
									metaInf.TotalFrags += 1
									if usePOS {
										metaInf.FragCounts[cdl4.F.Pos] += 1
									} else {
										metaInf.FragCounts[cdl4.F.Epos] += 1
									}
								}
								if cdl4.Cdl != nil {
									for _, cdlFive := range cdl4.Cdl {
										cdlChecker(cdlFive)
									}
								}
							}
							cdlChecker(cdlFour)
						}
					}
				}
			}
		}
	}

	// TODO TODO TODO: the returned text is being used as markdown, so there are issues with escape characters and such
	if printMetaInfo {
		metaInf.TotalLines = numLines
		metaInf.LineLengthAverage = float32(metaInf.TotalFrags) / float32(numLines)
		metaInf.TotalFragsWithoutBlanks = metaInf.TotalFrags - metaInf.FragCounts[""]
		metaInf.LineLengthAverageWithoutBlanks = float32(metaInf.TotalFragsWithoutBlanks) / float32(numLines)

		printMetaInfo := fmt.Sprintf("Meta information on text: %s\n\n"+
			"\tTotalFrags: %d\n"+
			"\tWithout Xs: %d\n"+
			"\tTotalLines: %d\n"+
			"\tLineLengthAverage: %.2f\n"+
			"\tAverageWithout Xs: %.2f\n",
			textBody.TextId,
			metaInf.TotalFrags,
			metaInf.TotalFragsWithoutBlanks,
			metaInf.TotalLines,
			metaInf.LineLengthAverage,
			metaInf.LineLengthAverageWithoutBlanks,
		)
		if metaInf.TotalFrags != 0 {
			metaFragCounts := fmt.Sprintln("\tFragCounts:")
			for s, i := range metaInf.FragCounts {
				metaFragCounts = fmt.Sprintf("%s\t\t%s\t%d\n", metaFragCounts, s, i)
			}
			printMetaInfo = fmt.Sprintf("%s%s", printMetaInfo, metaFragCounts)
		}
		textReturn = fmt.Sprintf("%s\n---\n%s", printMetaInfo, textReturn)
	}

	return fmt.Sprintf("%s\nEnd of text", textReturn)
}
