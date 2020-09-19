package inout

type TckrIn struct {
	Instrument string `json:"instrument_name"`
}

type TckrStats struct {
	Volume float64 `json:"volume"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
}

type TckrOut struct {
	BestAskAmt     float64   `json:"best_ask_amount"`
	BestAsk        float64   `json:"best_ask_price"`
	BestBidAmt     float64   `json:"best_bid_amount"`
	BestBid        float64   `json:"best_bid_price"`
	CurrentFunding float64   `json:"current_funding"`
	Funding8h      float64   `json:"funding_8h"`
	IndxPrc        float64   `json:"index_price"`
	Instrument     string    `json:"instrument_name"`
	Last           float64   `json:"last_price"`
	Mark           float64   `json:"mark_price"`
	MaxPrc         float64   `json:"max_price"`
	MinPrc         float64   `json:"min_price"`
	OpnIntrst      float64   `json:"open_interest"`
	StlmntPrc      float64   `json:"settlement_price"`
	State          string    `json:"state"`
	Stats          TckrStats `json:"stats"`
	TmStmp         int64     `json:"timestamp"`
}
