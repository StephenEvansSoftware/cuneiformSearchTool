package text

import (
	"fmt"
	"slices"
)

type TextMetaInfo struct {
	TotalFrags                     int
	TotalFragsWithoutBlanks        int
	FragCounts                     map[string]int
	TotalLines                     int
	LineLengthAverage              float32
	LineLengthAverageWithoutBlanks float32
	CfCounts                       map[string]int
	CfTypes                        map[string][]string
}

type PrintableText struct {
	TextId string
	Lines  []*Line

	MetaInformation *TextMetaInfo
}

type Line struct {
	Number int
	Label  string

	// This holds the frag / form and its cf
	Frags [][3]string
}

func (tmi *TextMetaInfo) Sprint() string {
	metaInfo := fmt.Sprintf("Meta information:\n\n"+
		"\tTotalFrags: %d\n"+
		"\tWithout Xs: %d\n"+
		"\tTotalLines: %d\n"+
		"\tLineLengthAverage: %.2f\n"+
		"\tAverageWithout Xs: %.2f\n",
		tmi.TotalFrags,
		tmi.TotalFragsWithoutBlanks,
		tmi.TotalLines,
		tmi.LineLengthAverage,
		tmi.LineLengthAverageWithoutBlanks,
	)
	if tmi.TotalFrags != 0 {
		metaInfo = fmt.Sprintf("%s%s", metaInfo, "\tFragCounts:\n")
		fragCountSlice := []string{}
		for s, i := range tmi.FragCounts {
			fragCountSlice = append(fragCountSlice, fmt.Sprintf("\t\t%s\t%d\n", s, i))
		}
		slices.Sort(fragCountSlice)
		for _, fcs := range fragCountSlice {
			metaInfo = fmt.Sprintf("%s%s", metaInfo, fcs)
		}
		metaInfo = fmt.Sprintf("%s%s", metaInfo, "\tCfCounts:\n")
		cfCountSlice := []string{}
		for s, i := range tmi.CfCounts {
			cfCountSlice = append(cfCountSlice, fmt.Sprintf("\t\t%s\t%d\n", s, i))
		}
		slices.Sort(cfCountSlice)
		for _, cfs := range cfCountSlice {
			metaInfo = fmt.Sprintf("%s%s", metaInfo, cfs)
		}
	}
	return fmt.Sprintf("%s\n---\n", metaInfo)
}

func (pt *PrintableText) Sprint() string {
	textReturn := fmt.Sprintf("Start of text: %s", pt.TextId)
	for _, line := range pt.Lines {
		textReturn = fmt.Sprintf("%s\n\t%d: %s\n\t\t", textReturn, line.Number, line.Label)
		for _, f := range line.Frags {
			textReturn = fmt.Sprintf("%s%s ", textReturn, f[0])
		}
	}
	textReturn = fmt.Sprintf("%s\nEnd of text", textReturn)
	return textReturn
}

// Note: highlights backwards
func (pt *PrintableText) Highlight(lineNumber int, fragIndex int, length int) string {
	//fmt.Println(pt.Lines[lineNumber-1].Frags[fragIndex-1])

	/*
		start one line X position Y
		go backwards on line X until we're done or hit position 0
		go backwards on line X-1 until we're done or hit position 0
	*/

	textReturn := fmt.Sprintf("Start of text: %s\n", pt.TextId)
	for ln, line := range pt.Lines {
		//textReturn = fmt.Sprintf("%s\n\t%d: %s\n\t\t", textReturn, line.Number, line.Label)
		textReturn = fmt.Sprintf("%s\n\t%d: %s\n\t\t", textReturn, line.Number, line.Label)
		for fn, f := range line.Frags {
			textReturn = fmt.Sprintf("%s*%s* ", textReturn, f[0])
			if ln == lineNumber-1 && fn == fragIndex-1 {
				//				textReturn = fmt.Sprintf("***%s***", textReturn)
			}
		}
	}
	textReturn = fmt.Sprintf("%s\nEnd of text", textReturn)
	return textReturn

}
