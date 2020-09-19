package inout

type CancelOrderIn struct {
	OrderID string `json:"order_id"`
}

type CancelAllByInstrmtIn struct {
	Instrument string `json:"instrument_name"`
	Type       string `json:"type,omitempty"`
}

type CancelAllByCcyIn struct {
	Ccy  string `json:"currency"`
	Kind string `json:"kind,omitempty"`
	Type string `json:"type,omitempty"`
}
