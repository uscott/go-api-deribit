package inout

type InstrumentIn struct {
	Ccy     string `json:"currency"`
	Kind    string `json:"kind,omitempty"`
	Expired bool   `json:"expired,omitempty"`
}

// InstrumentOut is the output argument for instrument queries
type InstrumentOut struct {
	BaseCcy         string  `json:"base_currency"`
	ContractSz      float64 `json:"contract_size"`
	CreatnTmStmp    int64   `json:"creation_timestamp"`
	ExprtnTmStmp    int64   `json:"expiration_timestamp"`
	Kind            string  `json:"kind"`
	Instrument      string  `json:"instrument_name"`
	IsActive        bool    `json:"is_active"`
	MakerCommission float64 `json:"maker_commission"`
	MinTrdAmt       float64 `json:"min_trade_amount"`
	OptnType        string  `json:"option_type"`
	QutCcy          string  `json:"quote_currency"`
	StlmntPrd       string  `json:"settlement_period"`
	Strike          float64 `json:"strike"`
	TakerCommission float64 `json:"taker_commission"`
	TckSz           float64 `json:"tick_size"`
}
