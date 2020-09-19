package inout

type CntrctSzIn struct {
	Instrument string `json:"instrument_name"`
}

type CntrctSzOut struct {
	CntrctSz float64 `json:"contract_size"`
}
