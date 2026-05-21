package text

// TODO: rename some of these structs

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

type Text struct {
	TextId string `json:"textid,omitempty"`

	Cdl []*CdlLevelOne `json:"cdl,omitempty"`
}

type CdlLevelOne struct {
	Cdl []*CdlLevelTwo `json:"cdl,omitempty"`
}

type CdlLevelTwo struct {
	Subtype string `json:"subtype,omitempty"`

	Cdl []*CdlLevelThree `json:"cdl,omitempty"`
}

type CdlLevelThree struct {
	Cdl []*CdlLevelFour `json:"cdl,omitempty"`
}

type CdlLevelFour struct {
	Frag  string `json:"frag,omitempty"`
	Label string `json:"label,omitempty"`
	Type  string `json:"type,omitempty"`

	F   *FTop           `json:"f,omitempty"`
	Cdl []*CdlLevelFour `json:"cdl,omitempty"`
}

type FTop struct {
	Cf   string `json:"cf,omitempty"`
	Epos string `json:"epos,omitempty"`
	Form string `json:"form,omitempty"`
	Pos  string `json:"pos,omitempty"`
}
