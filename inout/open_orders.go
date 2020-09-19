package inout

type OpnOrdrsByCcyIn struct {
	Ccy  string `json:"currency"`
	Kind string `json:"kind,omitempty"`
	Type string `json:"type,omitempty"`
}

type OpnOrdrsByInstrmtIn struct {
	Instrument string `json:"instrument_name"`
	Type       string `json:"type,omitempty"`
	Label      string `json:"label,omitempy"`
}
