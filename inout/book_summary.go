package inout

type BkSummaryByInstrmtIn struct {
	Instrument string `json:"instrument_name"`
}

type BkSummaryByCcyIn struct {
	Ccy  string `json:"currency"`
	Kind string `json:"kind,omitempty"`
}

type BkSummaryOut struct {
	Volume       float64 `json:"volume"`
	UndrlyngPrc  float64 `json:"underlying_price"`
	UndrlyngIndx string  `json:"underlying_index"`
	QutCcy       string  `json:"quote_currency"`
	OpnIntrst    float64 `json:"open_interest"`
	Mid          float64 `json:"mid_price"`
	Mark         float64 `json:"mark_price"`
	Low          float64 `json:"low"`
	Last         float64 `json:"last"`
	IntRate      float64 `json:"interest_rate"`
	Instrument   string  `json:"instrument_name"`
	High         float64 `json:"high"`
	CreatnTmStmp int64   `json:"creation_timestamp"`
	Bid          float64 `json:"bid_price"`
	BaseCcy      string  `json:"base_currency"`
	Ask          float64 `json:"ask_price"`
}
