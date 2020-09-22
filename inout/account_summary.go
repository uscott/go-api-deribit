package inout

// AcctSummaryIn is arguments to query for account summary
type AcctSummaryIn struct {
	Ccy      string `json:"currency"`
	Extended bool   `json:"extended,omitempty"`
}

type Rates struct {
	Burst int `json:"burst"`
	Rate  int `json:"rate"`
}

type RateLimits struct {
	Futures           Rates `json:"futures,omitempty"`
	MatchingEngine    Rates `json:"matching_engine,omitempty"`
	NonMatchingEngine Rates `json:"non_matching_engine,omitempty"`
	Options           Rates `json:"options,omitempty"`
	Perpetuals        Rates `json:"perpetuals,omitempty"`
}

// AcctSummaryOut contains info about your account
type AcctSummaryOut struct {
	AvailableFunds    float64    `json:"available_funds"`
	AvailableWdrFunds float64    `json:"available_withdrawal_funds"`
	Blnc              float64    `json:"balance"`
	Ccy               string     `json:"currency"`
	DeltaTotal        float64    `json:"delta_total"`
	DepositAddr       string     `json:"deposit_address"`
	Email             string     `json:"email"`
	Equity            float64    `json:"equity"`
	FuturesPl         float64    `json:"futures_pl"`
	FuturesSessionRpl float64    `json:"futures_session_rpl"`
	FuturesSessionUpl float64    `json:"futures_session_upl"`
	ID                int        `json:"id"`
	InitMrgn          float64    `json:"initial_margin"`
	Limits            RateLimits `json:"limits"`
	MntcMrgn          float64    `json:"maintenance_margin"`
	MrgnBlnc          float64    `json:"margin_balance"`
	OptnsDelta        float64    `json:"options_delta"`
	OptnsGamma        float64    `json:"options_gamma"`
	OptnsPl           float64    `json:"options_pl"`
	OptnsSessionRpl   float64    `json:"options_session_rpl"`
	OptnsSessionUpl   float64    `json:"options_session_upl"`
	OptnsTheta        float64    `json:"options_theta"`
	OptnsVega         float64    `json:"options_vega"`
	PrtflioMrgnEnbld  bool       `json:"portfolio_margining_enabled"`
	SessionFunding    float64    `json:"session_funding"`
	SessionRpl        float64    `json:"session_rpl"`
	SessionUpl        float64    `json:"session_upl"`
	SysName           string     `json:"system_name"`
	TfaEnbld          bool       `json:"tfa_enabled"`
	TotalPl           float64    `json:"total_pl"`
	Type              string     `json:"type"`
	UserName          string     `json:"username"`
}
