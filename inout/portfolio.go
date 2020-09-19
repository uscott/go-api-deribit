package inout

// Portfolio contains portolio information
type Portfolio struct {
	AvailableFunds    int    `json:"available_funds"`
	AvailableWdrFunds int    `json:"available_withdrawal_funds"`
	Bal               int    `json:"balance"`
	Ccy               string `json:"currency"`
	Equity            int    `json:"equity"`
	InitMrgn          int    `json:"initial_margin"`
	MntcMrgn          int    `json:"maintenance_margin"`
	MrgnBal           int    `json:"margin_balance"`
}
