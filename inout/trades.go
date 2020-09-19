package inout

type Trade struct {
	Amt         float64 `json:"amount"`
	Drctn       string  `json:"direction"`
	IndxPrc     float64 `json:"index_price"`
	Instrument  string  `json:"instrument_name"`
	IV          float64 `json:"iv"`
	Liquidation string  `json:"liquidation"`
	Prc         float64 `json:"price"`
	TckDrctn    int     `json:"tick_direction"`
	TmStmp      int64   `json:"timestamp"`
	TradeID     string  `json:"trade_id"`
	TrdSqnc     int     `json:"trade_seq"`
}

type UserTrade struct {
	Amt        float64     `json:"amount"`
	Drctn      string      `json:"direction"`
	Fee        float64     `json:"fee"`
	FeeCcy     string      `json:"fee_currency"`
	IndxPrc    float64     `json:"index_price"`
	Instrument string      `json:"instrument_name"`
	IV         string      `json:"iv"`
	Label      string      `json:"label"`
	Liquidity  string      `json:"liquidity"`
	MatchingID interface{} `json:"matching_id"`
	OrderID    string      `json:"order_id"`
	OrderType  string      `json:"order_type"`
	PostOnly   bool        `json:"post_only"`
	Prc        float64     `json:"price"`
	RdcOnly    bool        `json:"reduce_only"`
	SelfTrd    bool        `json:"self_trade"`
	State      string      `json:"state"`
	TckDrctn   int         `json:"tick_direction"`
	TmStmp     int64       `json:"timestamp"`
	TradeID    string      `json:"trade_id"`
	TrdSqnc    int         `json:"trade_seq"`
	UndPrc     float64     `json:"underlying_price"`
}

type UserTradesOut struct {
	HasMore bool        `json:"has_more"`
	Trades  []UserTrade `json:"trades"`
}

type LastTradesOut struct {
	HasMore bool    `json:"has_more"`
	Trades  []Trade `json:"trades"`
}

type TradesByInstrmtIn struct {
	Instrument string `json:"instrument_name"`
	StartSeq   int64  `json:"start_seq,omitempty"`
	EndSeq     int64  `json:"end_seq,omitempty"`
	Count      int    `json:"count,omitempty"`
	IncludeOld bool   `json:"include_old,omitempty"`
	Sorting    string `json:"sorting,omitempty"`
}

type TradesByInstrmtAndTmIn struct {
	Instrument  string `json:"instrument_name"`
	StartTmStmp int64  `json:"start_timestamp"`
	EndTmStmp   int64  `json:"end_timestamp"`
	Count       int    `json:"count,omitempty"`
	IncludeOld  bool   `json:"include_old,omitempty"`
	Sorting     string `json:"sorting,omitempty"`
}

type TradesByCcyIn struct {
	Ccy        string `json:"currency"`
	Kind       string `json:"kind,omitempty"`
	StartId    int64  `json:"start_id,omitempty"`
	EndId      int64  `json:"end_id,omitempty"`
	Count      int    `json:"count,omitempty"`
	IncludeOld bool   `json:"include_old,omitempty"`
	Sorting    string `json:"sorting,omitempty"`
}

type TradesByCcyAndTmIn struct {
	Ccy         string `json:"currency"`
	Kind        string `json:"kind,omitempty"`
	StartTmStmp int64  `json:"start_timestamp"`
	EndTmStmp   int64  `json:"end_timestamp"`
	Count       int    `json:"count,omitempty"`
	IncludeOld  bool   `json:"include_old,omitempty"`
	Sorting     string `json:"sorting,omitempty"`
}
