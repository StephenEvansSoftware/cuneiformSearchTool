package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"slices"

	"cuneiformSearchTool/text"
)

var (
	// Valid word types
	ValidTypes = []string{
		"AJ",  // Adjective
		"AV",  // Adverb
		"CN",  // Celestial name
		"CNJ", // Conjunction
		"DET", // Determinative pronoun
		"DN",  // Divine name
		"DP",  // Demonstrative pronoun
		"EN",  // Ethnic name
		"GN",  // Geographic name
		"IP",  // Independant pronoun
		"J",   // Exclamation / interjection
		"MN",  // Month name
		"MOD", // Modifier
		"N",   // Noun
		"NU",  // Numeral
		"ON",  // Settlement name
		"PN",  // Personal name
		"PP",  // Possesive pronoun
		"PRP", // Preposition
		"QN",  // Place name
		"QP",  // Question particle
		"REL", // Relative pronoun
		"RN",  // Ruler name
		"RP",  // Reflexive prnoun
		"SBJ", // Subjunction
		"SN",  // City / settlement name
		"TN",  // Temple name
		"V",   // Verb
		"V/i", //// Some sumerian stuff
		"V/t", //// Some sumerian stuff
		"WN",  // Water name
		"X",   // Meaning unknown
		"XP",  // Indefinite pronoun
		"n",   //// Breaks? Only on POS
		"u",   //// Breaks? Only on POS
	}

	AllProperNounsType = "All proper nouns"

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

	validTypesSansProperNouns = []string{
		"AJ",
		"AV",
		"CNJ",
		"DET",
		"DP",
		"IP",
		"J",
		"MOD",
		"N",
		"NU",
		"PP",
		"PRP",
		"QP",
		"REL",
		"RP",
		"SBJ",
		"V",
		"V/i",
		"V/t",
		"X",
		"XP",
		"n",
		"u",
	}

	inputCorpus  = "./jsonFiles/corpus"

	corpusTextInfo = []*CorpusInfo{}

	totalTermsWithoutProperNounGrouping int
	totalTermsWithProperNounGrouping    int
)

type CorpusSearchResult struct {
	CorpusName  string
	TextResults []*TextSearchResult
}

type TextSearchResult struct {
	TextName string
	Results  []*SearchResult
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
}

type MassSearchSettings struct {
	// Terms length we want to mass search for
	TermsLength int

	// Do we want to group proper nouns for this mass search
	UseProperNounGrouping bool

	// Do we want to do >13k searches all at once?
	LimitSearchScope bool
	// Search scope of number of terms
	SearchFrom int
	SearchTo   int

	// Can manually set a selection of terms
	ManualTermSelection bool
	// TODO: make a struct for this. 3d slice is annoying to work with
	ManualTermsToSearch [][][]string
}

type PrintSettings struct {
	PrintTextMetaInformation bool
	UsePOS                   bool
	UseFORM                  bool
}

type CorpusMetaInfo struct {
	NumberOfTexts                  int
	TotalFragsInAllTexts           int
	TotalFragsInAllTextsWithoutX   int
	FragCounts                     map[string]int
	TotalLinesInAllTexts           int
	AverageLineLength              float32
	AverageLineLengthWithoutBlanks float32
}

// Package intiialisation
func init() {
	if err := loadCorpusTexts(); err != nil {
		panic(err.Error())
	}

	totalTermsWithoutProperNounGrouping = len(ValidTypes)
	totalTermsWithProperNounGrouping = len(validTypesSansProperNouns) + 1
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
			// Special case where the proper nouns group is the term
			if search == AllProperNounsType {
				break
			}
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

func getCorpusList(allCorpus bool, selectCorpus []string) []string {
	// When allCorpus is set, change which corpus to search by overriding its value
	if allCorpus {
		selectCorpus = []string{}
		for _, corp := range corpusTextInfo {
			selectCorpus = append(selectCorpus, corp.CorpusName)
		}
	}
	return selectCorpus
}

// Main search functionality
func SearchFor(terms [][]string, searchSettings *SearchSettings) ([]*CorpusSearchResult, int, error) {
	searchSettings.SelectedCorpus = getCorpusList(searchSettings.SearchAllTexts, searchSettings.SelectedCorpus)

	if len(searchSettings.SelectedCorpus) == 0 {
		// No corpus' selected, return error
		return nil, -1, fmt.Errorf("no corpus selected")
	}

	// Results and total count of terms within corpus
	corpusSearchResults := []*CorpusSearchResult{}
	totalInstancesFound := 0

	for _, corpTextInfo := range corpusTextInfo {
		for _, corp := range searchSettings.SelectedCorpus {
			if corp == corpTextInfo.CorpusName {
				// Search this corpus
				textSearchResults := []*TextSearchResult{}
				// Range over each text within the corpus
				for _, textInfo := range corpTextInfo.Texts {
					// TODO TODO TODO: add text deliniation here
					textBytes, err := readFile(textInfo.TextPath)
					if err != nil {
						return nil, -1, err
					}
					text, err := marshalToText(textBytes)
					if err != nil {
						return nil, -1, err
					}
					// Perform the actual search
					searchResult, err := searchText(text, terms, searchSettings)
					if err != nil {
						return nil, -1, err
					}
					// Avoid adding empty result to search results
					if len(searchResult.Results) != 0 {
						searchResult.TextPath = textInfo.TextPath
						textSearchResults = append(textSearchResults, searchResult)
						totalInstancesFound += len(searchResult.Results)
					}
				}
				// TODO: make this a setting to hide empty results or not
				// if len(textSearchResults) != 0 {
				corpusSearchResults = append(corpusSearchResults, &CorpusSearchResult{
					CorpusName:  corp,
					TextResults: textSearchResults,
				})
				// }
			}
		}
	}
	return corpusSearchResults, totalInstancesFound, nil
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

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileBytes, _ := ioutil.ReadAll(file)
	return fileBytes, nil
}

func marshalToText(bytes []byte) (*text.Text, error) {
	var text text.Text
	if err := json.Unmarshal(bytes, &text); err != nil {
		return nil, err
	}
	return &text, nil
}

// TODO: a lot. this was just taken pretty much verbatim from my initial prototype god-file monstrosity
func searchText(textBody *text.Text, searchPattern [][]string, searchSettings *SearchSettings) (*TextSearchResult, error) {
	// Ensure we're actually searching for something
	if len(searchPattern) == 0 {
		return nil, fmt.Errorf("cannot have empty search paattern")
	}

	// Check if "All proper nouns" are being searched for
	for i, sp := range searchPattern {
		for j, spt := range sp {
			if spt == AllProperNounsType {
				// Overwrite the "All proper nouns" text with the first proper noun
				searchPattern[i][j] = properNounTypes[0]
				// Append the remaining proper nouns from the list to the search pattern
				searchPattern[i] = append(searchPattern[i], properNounTypes[1:]...)
			}
		}
	}

	spLen := len(searchPattern)
	// Slice of bools that denote if we have found up to the nth position of the terms so far
	previouslyFoundPatternPart := make([]bool, spLen)

	// Collection of results
	foundPatterns := &TextSearchResult{
		TextName: textBody.TextId,
		Results:  []*SearchResult{},
	}

	// Index in the current line of the current fragment
	fragIndex := 0
	// Line count
	numLines := 1
	// Line label
	currLabel := ""

	// TODO TODO TODO: messy and gross looping through so many structs even if they are the original structure of the data
	// TODO TODO TODO: change the structure to something easier to manage
	for _, cdlOne := range textBody.Cdl {
		for _, cdlTwo := range cdlOne.Cdl {
			if cdlTwo.Subtype == "body" {
				// Here we have the actual body of the text
				for _, cdlThree := range cdlTwo.Cdl {
					if len(cdlThree.Cdl) != 0 {
						for _, cdlFour := range cdlThree.Cdl {
							// Since there can a cdlFour inside of a cdlFour, recursively check it
							var cdlChecker func(cdl4 *text.CdlLevelFour)
							cdlChecker = func(cdl4 *text.CdlLevelFour) {
								// Start of the actual text analysis
								if cdl4.Type == "line-start" {
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
									// This is the actual element in the text we care about
									fragIndex += 1
									// Reverse iterate through from element N-1 to 0
									// pfppi = previousFoundPatternPartIndex
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
												break
											}
										} else {
											if cdl4.F.Epos == spf {
												previouslyFoundPatternPart[0] = true
												break
											}
										}
									}
									// Final check to see if we have completed the pattern or not
									if previouslyFoundPatternPart[spLen-1] {
										// WE HAVE FOUND THE PATTERN
										// Store line and number
										newResult := &SearchResult{
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
										// Reset last index in case the pattern was found
										previouslyFoundPatternPart[spLen-1] = false
									}
								}
								if cdl4.Cdl != nil {
									// Recursively check cdl within current cdl
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

	return foundPatterns, nil
}

func GetPrintableText(textPath string, printSettings *PrintSettings) (string, *text.TextMetaInfo, error) {
	textBytes, err := readFile(textPath)
	if err != nil {
		return "", nil, err
	}
	text, err := marshalToText(textBytes)
	if err != nil {
		return "", nil, err
	}
	printText, textMetaInfo := getPrintableText(text, printSettings)
	return printText, textMetaInfo, nil
}

func getPrintableText(textBody *text.Text, printSettings *PrintSettings) (string, *text.TextMetaInfo) {
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
									if printSettings.UseFORM {
										textReturn = fmt.Sprintf("%s%s ", textReturn, cdl4.F.Form)
									} else {
										textReturn = fmt.Sprintf("%s%s ", textReturn, cdl4.Frag)
									}
									metaInf.TotalFrags += 1
									if printSettings.UsePOS {
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
	if printSettings.PrintTextMetaInformation {
		metaInf.TotalLines = numLines - 1
		metaInf.LineLengthAverage = float32(metaInf.TotalFrags) / float32(numLines-1)
		metaInf.TotalFragsWithoutBlanks = metaInf.TotalFrags - metaInf.FragCounts[""]
		metaInf.LineLengthAverageWithoutBlanks = float32(metaInf.TotalFragsWithoutBlanks) / float32(numLines-1)

		metaInfo := fmt.Sprintf("Meta information on text: %s\n\n"+
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
			metaInfo = fmt.Sprintf("%s%s", metaInfo, "\tFragCounts:\n")
			fragCountSlice := []string{}
			for s, i := range metaInf.FragCounts {
				fragCountSlice = append(fragCountSlice, fmt.Sprintf("\t\t%s\t%d\n", s, i))
			}
			slices.Sort(fragCountSlice)
			for _, fcs := range fragCountSlice {
				metaInfo = fmt.Sprintf("%s%s", metaInfo, fcs)
			}
		}
		textReturn = fmt.Sprintf("%s\n---\n%s", metaInfo, textReturn)
	}

	return fmt.Sprintf("%s\nEnd of text", textReturn), metaInf
}

func GetCorpusMetaInfo(corpus string) (*CorpusMetaInfo, error) {
	for _, c := range corpusTextInfo {
		if c.CorpusName == corpus {
			cmi := &CorpusMetaInfo{FragCounts: map[string]int{}}
			cmi.NumberOfTexts = len(c.Texts)
			printTextSettings := &PrintSettings{PrintTextMetaInformation: true}
			for _, t := range c.Texts {
				_, tmi, err := GetPrintableText(t.TextPath, printTextSettings)
				if err != nil {
					panic(err.Error()) // TODO: handle better
				}
				cmi.TotalFragsInAllTexts += tmi.TotalFrags
				cmi.TotalFragsInAllTextsWithoutX += tmi.TotalFragsWithoutBlanks
				cmi.TotalLinesInAllTexts += tmi.TotalLines

				for s, i := range tmi.FragCounts {
					cmi.FragCounts[s] += i
				}
			}

			cmi.AverageLineLength = float32(cmi.TotalFragsInAllTexts) / float32(cmi.TotalLinesInAllTexts)
			cmi.AverageLineLengthWithoutBlanks = float32(cmi.TotalFragsInAllTextsWithoutX) / float32(cmi.TotalLinesInAllTexts)

			return cmi, nil
		}
	}
	return nil, fmt.Errorf("corpus '%s' not found", corpus)
}
