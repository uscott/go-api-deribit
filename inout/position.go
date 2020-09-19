package inout

type PosnInstrmtIn struct {
	Instrument string `json:"instrument_name"`
}

type PosnCcyIn struct {
	Ccy  string `json:"currency"`
	Kind string `json:"kind,omitempty"`
}

type PosnOut struct {
	AvgPrc       float64 `json:"average_price"`
	Delta        float64 `json:"delta"`
	Drctn        string  `json:"direction"`
	EstdLqdtnPrc float64 `json:"estimated_liquidation_price"`
	FltngPl      float64 `json:"floating_profit_loss"`
	IndxPrc      float64 `json:"index_price"`
	InitMrgn     float64 `json:"initial_margin"`
	Instrument   string  `json:"instrument_name"`
	Kind         string  `json:"kind"`
	Leverage     float64 `json:"leverage"`
	MaintMrgn    float64 `json:"maintenance_margin"`
	Mark         float64 `json:"mark_price"`
	OpnOrdrsMrgn float64 `json:"open_orders_margin"`
	RlzdPl       float64 `json:"realized_profit_loss"`
	StlmntPrc    float64 `json:"settlement_price"`
	Sz           float64 `json:"size"`
	SzCcy        float64 `json:"size_currency"`
	TotalPl      float64 `json:"total_profit_loss"`
}
