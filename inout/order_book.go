package inout

type BookIn struct {
	Instrument string `json:"instrument_name"`
	Depth      int    `json:"depth,omitempty"`
}

type BookOut struct {
	TmStmp         int64       `json:"timestamp"`
	Stats          TckrStats   `json:"stats"`
	State          string      `json:"state"`
	StlmntPrc      float64     `json:"settlement_price"`
	OpnIntrst      float64     `json:"open_interest"`
	MinPrc         float64     `json:"min_price"`
	MaxPrc         float64     `json:"max_price"`
	Mark           float64     `json:"mark_price"`
	Last           float64     `json:"last_price"`
	Instrument     string      `json:"instrument_name"`
	IndxPrc        float64     `json:"index_price"`
	Funding8h      float64     `json:"funding_8h"`
	CurrentFunding float64     `json:"current_funding"`
	ChgId          int         `json:"change_id"`
	Bids           [][]float64 `json:"bids"`
	BestBid        float64     `json:"best_bid_price"`
	BestBidAmt     float64     `json:"best_bid_amount"`
	BestAsk        float64     `json:"best_ask_price"`
	BestAskAmt     float64     `json:"best_ask_amount"`
	Asks           [][]float64 `json:"asks"`
}
