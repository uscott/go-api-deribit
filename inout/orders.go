package inout

import "strconv"

type OrderIn struct {
	Instrument string   `json:"instrument_name"`
	Amt        float64  `json:"amount"`
	Type       string   `json:"type,omitempty"`
	Label      string   `json:"label,omitempty"`
	Prc        float64  `json:"price,omitempty"`
	TmInFrc    string   `json:"time_in_force,omitempty"`
	MaxShow    *float64 `json:"max_show,omitempty"`
	PostOnly   bool     `json:"post_only,omitempty"`
	RdcOnly    bool     `json:"reduce_only,omitempty"`
	StopPrc    float64  `json:"stop_price,omitempty"`
	Trigger    string   `json:"trigger,omitempty"`
	Advanced   string   `json:"advanced,omitempty"`
}

type EditIn struct {
	OrderID  string  `json:"order_id"`
	Amt      float64 `json:"amount"`
	Prc      float64 `json:"price"`
	PostOnly bool    `json:"post_only,omitempty"`
	Advanced string  `json:"advanced,omitempty"`
	StopPrc  float64 `json:"stop_price,omitempty"`
}

type OrderOut struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}

type Price float64

func (p *Price) UnmarshalJSON(data []byte) (err error) {
	if string(data) == `"market_price"` {
		*p = 0
		return
	}
	var f float64
	f, err = strconv.ParseFloat(string(data), 0)
	if err != nil {
		return
	}
	*p = Price(f)
	return
}

func (p *Price) ToFloat64() float64 {
	return float64(*p)
}

// Order contains the information for a submitted order
type Order struct {
	Advanced      string  `json:"advanced,omitempty"`
	Amt           float64 `json:"amount"`
	API           bool    `json:"api"`
	AvgPrc        float64 `json:"average_price"`
	Commission    float64 `json:"commission"`
	CreatnTmStmp  int64   `json:"creation_timestamp"`
	Drctn         string  `json:"direction"`
	FilledAmt     float64 `json:"filled_amount"`
	Implv         float64 `json:"implv,omitempty"`
	Instrument    string  `json:"instrument_name"`
	IsLiquidation bool    `json:"is_liquidation"`
	Label         string  `json:"label"`
	LstUpdtTmStmp int64   `json:"last_update_timestamp"`
	OrderID       string  `json:"order_id"`
	OrderState    string  `json:"order_state"`
	OrderType     string  `json:"order_type"`
	MaxShow       float64 `json:"max_show"`
	PostOnly      bool    `json:"post_only"`
	Prc           Price   `json:"price"`
	ProfitLoss    float64 `json:"profit_loss"`
	RdcOnly       bool    `json:"reduce_only"`
	TmInFrc       string  `json:"time_in_force"`
	StopPrc       float64 `json:"stop_price,omitempty"`
	Triggered     bool    `json:"triggered,omitempty"`
	USD           float64 `json:"usd,omitempty"`
}

// StopOrderHistoryIn is the input parameter for GetStopOrderHistory
type StopOrderHistoryIn struct {
	Ccy          string `json:"currency"`
	Instrument   string `json:"instrument_name,omitempty"`
	Count        int    `json:"count,omitempty"`
	Continuation string `json:"continuation,omitempty"`
}

// StopOrder contains info for a submitted stop order
type StopOrder struct {
	Trigger    string  `json:"trigger"`
	TmStmp     int64   `json:"timestamp"`
	StopPrice  float64 `json:"stop_price"`
	StopID     string  `json:"stop_id"`
	OrderState string  `json:"order_state"`
	Request    string  `json:"request"`
	Prc        float64 `json:"price"`
	OrderID    string  `json:"order_id"`
	Offset     int     `json:"offset"`
	Instrument string  `json:"instrument_name"`
	Amt        float64 `json:"amount"`
	Drctn      string  `json:"direction"`
}
