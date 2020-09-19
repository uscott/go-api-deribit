package inout

type IndxIn struct {
	Ccy string `json:"currency"`
}

type IndxOut struct {
	BTC float64 `json:"BTC"`
	ETH float64 `json:"ETH"`
	EDP float64 `json:"edp"`
}
