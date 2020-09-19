package inout

// XferSubAcctIn is the argument struct for the
// inputs to posting a request for a transfer
// between (sub) accounts
type XferSubAcctIn struct {
	Amt string `json:"amount"`
	Ccy string `json:"currency"`
	Dst int    `json:"destination"`
}

// XferSubAcctOut is the struct containing the
// result form a transfer request
type XferSubAcctOut struct {
	Amt           float64 `json:"amount"`
	CreatedTmStmp int64   `json:"created_timestamp"`
	Currency      string  `json:"currency"`
	Direction     string  `json:"direction"`
	ID            int     `json:"id"`
	OtherSide     string  `json:"other_side"`
	State         string  `json:"state"`
	Type          string  `json:"type"`
	UpdatedTmStmp int64   `json:"updated_timestamp"`
}
