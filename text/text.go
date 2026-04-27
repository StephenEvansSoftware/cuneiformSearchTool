package text

// TODO: rename some of these structs

type TextMetaInfo struct {
	TotalFrags                     int
	TotalFragsWithoutBlanks        int
	FragCounts                     map[string]int
	TotalLines                     int
	LineLengthAverage              float32
	LineLengthAverageWithoutBlanks float32
}

type Text struct {
	License    string `json:"license,omitempty"`
	LicenceUrl string `json:"license-url,omitempty"`
	MoreInfo   string `json:"more-info,omitempty"`
	Project    string `json:"project,omitempty"`
	Source     string `json:"source,omitempty"`
	TextId     string `json:"textid,omitempty"`
	Timestamp  string `json:"UTC-timestamp,omitempty"`
	Type       string `json:"type,omitempty"`

	Cdl []*CdlLevelOne `json:"cdl,omitempty"`
}

type CdlLevelOne struct {
	Id   string `json:"id,omitempty"`
	Node string `json:"node,omitempty"`
	Type string `json:"type,omitempty"`

	Cdl      []*CdlLevelTwo `json:"cdl,omitempty"`
	Linkbase []*LinkbaseCdl `json:"linkbase,omitempty"`
}

type LinkbaseCdl struct {
	Id         string           `json:"id,omitempty"`
	Sig        string           `json:"sig,omitempty"`
	Type       string           `json:"type,omitempty"`
	XlinkRole  string           `json:"xlink_role,omitempty"`
	XlinkTitle string           `json:"xlink_title,omitempty"`
	XlinkType  string           `json:"xlink_type,omitempty"`
	Links      []*LinksFromBase `json:"links,omitempty"`
}

type LinksFromBase struct {
	Type       string `json:"type,omitempty"`
	XlinkHref  string `json:"xlink_href,omitempty"`
	XlinkRole  string `json:"xlink_role,omitempty"`
	XlinkTitle string `json:"xlink_title,omitempty"`
	XlinkType  string `json:"xlink_type,omitempty"`
}

type CdlLevelTwo struct {
	Extend  string `json:"extent,omitempty"`
	Id      string `json:"id,omitempty"`
	Label   string `json:"label,omitempty"`
	N       string `json:"n,omitempty"`
	Node    string `json:"node,omitempty"`
	Queried string `json:"queried,omitempty"`
	Ref     string `json:"ref,omitempty"`
	Scope   string `json:"scope,omitempty"`
	State   string `json:"state,omitempty"`
	Strict  string `json:"strict,omitempty"`
	Subtype string `json:"subtype,omitempty"`
	Type    string `json:"type,omitempty"`

	Cdl []*CdlLevelThree `json:"cdl,omitempty"`
}

type CdlLevelThree struct {
	Extent   string `json:"extent,omitempty"`
	Id       string `json:"id,omitempty"`
	Implicit string `json:"implicit,omitempty"`
	Label    string `json:"label,omitempty"`
	N        string `json:"n,omitempty"`
	Node     string `json:"node,omitempty"`
	Queried  string `json:"queried,omitempty"`
	Ref      string `json:"ref,omitempty"`
	Scope    string `json:"scope,omitempty"`
	State    string `json:"state,omitempty"`
	Strict   string `json:"strict,omitempty"`
	Subtype  string `json:"subtype,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Type     string `json:"type,omitempty"`

	Cdl []*CdlLevelFour `json:"cdl,omitempty"`
}

type CdlLevelFour struct {
	Bad           string `json:"bad,omitempty"`
	CofHead       string `json:"cof-head,omitempty"`
	CofTails      string `json:"cof-tails,omitempty"`
	Exolng        string `json:"exolng,omitempty"`
	Exoprj        string `json:"exoprj,omitempty"`
	Exosig        string `json:"exosig,omitempty"`
	Extent        string `json:"extent,omitempty"`
	Flags         string `json:"flags,omitempty"`
	Frag          string `json:"frag,omitempty"`
	Id            string `json:"id,omitempty"`
	Inst          string `json:"inst,omitempty"`
	Label         string `json:"label,omitempty"`
	Lang          string `json:"lang,omitempty"`
	N             string `json:"n,omitempty"`
	Newsig        string `json:"newsig,omitempty"`
	Node          string `json:"node,omitempty"`
	Queried       string `json:"queried,omitempty"`
	Ref           string `json:"ref,omitempty"`
	Scope         string `json:"scope,omitempty"`
	ScoreBlockRef string `json:"score-block-ref,omitempty"`
	Sig           string `json:"sig,omitempty"`
	State         string `json:"state,omitempty"`
	Strict        string `json:"strict,omitempty"`
	Subtype       string `json:"subtype,omitempty"`
	Tag           string `json:"tag,omitempty"`
	TailSig       string `json:"tail-sig,omitempty"`
	Type          string `json:"type,omitempty"`

	F     *FTop           `json:"f,omitempty"`
	Para  []*ParaTop      `json:"para,omitempty"`
	Props []*PropsTop     `json:"props,omitempty"`
	Cdl   []*CdlLevelFour `json:"cdl,omitempty"` // TODO: not sure if this will work
}

type ParaTop struct {
	BracketingLevel string `json:"bracketing_level,omitempty"`
	Class           string `json:"class,omitempty"`
	Text            string `json:"text,omitempty"`
	Type            string `json:"type,omitempty"`
}

type PropsTop struct {
	Group string `json:"group,omitempty"`
	Name  string `json:"name,omitempty"`
	Ngram string `json:"ngram,omitempty"`
	Ref   string `json:"ref,omitempty"`
	Value string `json:"value,omitempty"`
}

type FTop struct {
	Base     string `json:"base,omitempty"`
	Cf       string `json:"cf,omitempty"`
	Cont     string `json:"cont,omitempty"`
	Contrefs string `json:"contrefs,omitempty"`
	Delim    string `json:"delim,omitempty"`
	Epos     string `json:"epos,omitempty"`
	Form     string `json:"form,omitempty"`
	Gw       string `json:"gw,omitempty"`
	Headform string `json:"headform,omitempty"`
	Lang     string `json:"lang,omitempty"`
	Morph    string `json:"morph,omitempty"`
	Norm     string `json:"norm,omitempty"`
	Norm0    string `json:"norm0,omitempty"`
	Pos      string `json:"pos,omitempty"`
	Sense    string `json:"sense,omitempty"`

	Gdl []*GdlTop `json:"gdl,omitempty"`
}

type GdlTop struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	Delim       string `json:"delim,omitempty"`
	Det         string `json:"det,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlEm       string `json:"gdl_em,omitempty"`
	GdlRemarked string `json:"gdl_remarked,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	GdlType     string `json:"gdl_type,omitempty"`
	Gg          string `json:"gg,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	N           string `json:"n,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Pos         string `json:"pos,omitempty"`
	Queried     string `json:"queried,omitempty"`
	S           string `json:"s,omitempty"`
	Sexified    string `json:"sexified,omitempty"`
	Spform      string `json:"spform,omitempty"`
	Spoid       string `json:"spoid,omitempty"`
	StatusStart string `json:"statusStart,omitempty"`
	Surro       string `json:"surro,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Uflag3      string `json:"uflag3,omitempty"`
	Utf8        string `json:"utf8,omitempty"`
	V           string `json:"v,omitempty"`
	X           string `json:"x,omitempty"`

	Group []*GroupGdl `json:"group,omitempty"`
	Mods  []*ModsGdl  `json:"mods,omitempty"`
	Seq   []*SeqGdl   `json:"seq,omitempty"`
}

type GroupGdl struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	C           string `json:"c,omitempty"`
	Delim       string `json:"delim,omitempty"`
	Det         string `json:"det,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlEm       string `json:"gdl_em,omitempty"`
	GdlRemarked string `json:"gdl_remarked,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	GdlType     string `json:"gdl_type,omitempty"`
	Gg          string `json:"gg,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Logolang    string `json:"logolang,omitempty"`
	N           string `json:"n,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	P           string `json:"p,omitempty"`
	Pos         string `json:"pos,omitempty"`
	Queried     string `json:"queried,omitempty"`
	Role        string `json:"role,omitempty"`
	S           string `json:"s,omitempty"`
	Sexified    string `json:"sexified,omitempty"`
	SpForm      string `json:"spform,omitempty"`
	SpOid       string `json:"spoid,omitempty"`
	StatusStart string `json:"statusStart,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Uflag3      string `json:"uflag3,omitempty"`
	Utf8        string `json:"utf8,omitempty"`
	V           string `json:"v,omitempty"`
	X           string `json:"x,omitempty"`

	Group []*GroupGroupGdl `json:"group,omitempty"`
	Mods  []*ModsGroupGdl  `json:"mods,omitempty"`
	Seq   []*SeqGroupGdl   `json:"seq,omitempty"`
}

type GroupGroupGdl struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	Delim       string `json:"delim,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	GdlType     string `json:"gdl_type,omitempty"`
	Gg          string `json:"gg,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Logolang    string `json:"logolang,omitempty"`
	N           string `json:"n,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Queried     string `json:"queried,omitempty"`
	Role        string `json:"role,omitempty"`
	S           string `json:"s,omitempty"`
	Sexified    string `json:"sexified,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Utf8        string `json:"utf8,omitempty"`
	V           string `json:"v,omitempty"`

	Group []*GroupGroupGroupGdl `json:"group,omitempty"`
	Seq   []*SeqGroupGroupGdl   `json:"seq,omitempty"`
}

type SeqGroupGroupGdl struct {
	R string `json:"r,omitempty"`
}

type GroupGroupGroupGdl struct {
	Break       string `json:"break,omitempty"`
	Delim       string `json:"delim,omitempty"`
	Det         string `json:"det,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	GdlSurroEnd string `json:"gdl_surroEnd,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Logolang    string `json:"logolang,omitempty"`
	N           string `json:"n,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Pos         string `json:"pos,omitempty"`
	Queried     string `json:"queried,omitempty"`
	Role        string `json:"role,omitempty"`
	S           string `json:"s,omitempty"`
	Sexified    string `json:"sexified,omitempty"`
	Utf8        string `json:"utf8,omitempty"`

	Seq []*SeqGroupGroupGroupGdl `json:"seq,omitempty"`
}

type SeqGroupGroupGroupGdl struct {
	Id   string `json:"id,omitempty"`
	R    string `json:"r,omitempty"`
	Utf8 string `json:"utf8,omitempty"`
	V    string `json:"v,omitempty"`
}

type ModsGroupGdl struct {
	B string `json:"b,omitempty"`
	F string `json:"f,omitempty"`
	M string `json:"m,omitempty"`
}

type SeqGroupGdl struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	Delim       string `json:"delim,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	GdlType     string `json:"gdl_type,omitempty"`
	Gg          string `json:"gg,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Logolang    string `json:"logolang,omitempty"`
	N           string `json:"n,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Queried     string `json:"queried,omitempty"`
	R           string `json:"r,omitempty"`
	Role        string `json:"role,omitempty"`
	S           string `json:"s,omitempty"`
	Sexified    string `json:"sexified,omitempty"`
	StatusStart string `json:"statusStart,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Utf8        string `json:"utf8,omitempty"`
	V           string `json:"v,omitempty"`

	Group []*GroupSeqGroupGdl `json:"group,omitempty"`
	Mods  []*ModsSeqGroupGdl  `json:"mods,omitempty"`
	Seq   []*SeqSeqGroupGdl   `json:"seq,omitempty"`
}

type ModsSeqGroupGdl struct {
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
}

type GroupSeqGroupGdl struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Logolang    string `json:"logolang,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Queried     string `json:"queried,omitempty"`
	Role        string `json:"role,omitempty"`
	S           string `json:"s,omitempty"`
	StatusStart string `json:"statusStart,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Utf8        string `json:"utf8,omitempty"`

	Mods []*ModsGroupSeq `json:"mods,omitempty"`
}

type SeqSeqGroupGdl struct {
	R string `json:"r,omitempty"`
}

type ModsGdl struct {
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
	F string `json:"f,omitempty"`
}

type SeqGdl struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	GdlType     string `json:"gdl_type,omitempty"`
	Gg          string `json:"gg,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	N           string `json:"n,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Queried     string `json:"queried,omitempty"`
	R           string `json:"r,omitempty"`
	S           string `json:"s,omitempty"`
	Sexified    string `json:"sexified,omitempty"`
	StatusStart string `json:"statusStart,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Utf8        string `json:"utf8,omitempty"`
	V           string `json:"v,omitempty"`

	Seq   []*SeqSeq      `json:"seq,omitempty"`
	Group []*GroupSeqGdl `json:"group,omitempty"`
}

type SeqSeq struct {
	R string `json:"r,omitempty"`
}

type GroupSeqGdl struct {
	Break       string `json:"break,omitempty"`
	BreakEnd    string `json:"breakEnd,omitempty"`
	BreakStart  string `json:"breakStart,omitempty"`
	Form        string `json:"form,omitempty"`
	GdlCollated string `json:"gdl_collated,omitempty"`
	GdlSign     string `json:"gdl_sign,omitempty"`
	Hc          string `json:"hc,omitempty"`
	Ho          string `json:"ho,omitempty"`
	Id          string `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	Logolang    string `json:"logolang,omitempty"`
	O           string `json:"o,omitempty"`
	Oid         string `json:"oid,omitempty"`
	Queried     string `json:"queried,omitempty"`
	Role        string `json:"role,omitempty"`
	S           string `json:"s,omitempty"`
	StatusStart string `json:"statusStart,omitempty"`
	Uflag1      string `json:"uflag1,omitempty"`
	Uflag3      string `json:"uflag3,omitempty"`
	Utf8        string `json:"utf8,omitempty"`

	Mods []*ModsGroupSeq `json:"mods,omitempty"`
}

type ModsGroupSeq struct {
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
	F string `json:"f,omitempty"`
}
