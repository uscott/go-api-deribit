package inout

import (
	"fmt"
	"strconv"
	"strings"
)

type NtfctnAnncmnt struct {
	Actn      string `json:"action"`
	Body      string `json:"body"`
	Date      int64  `json:"date"`
	ID        int    `json:"id"`
	Important bool   `json:"important"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
}

type NtfctnDrbtPrcIndx struct {
	TmStmp int64   `json:"timestamp"`
	Prc    float64 `json:"price"`
	Indx   string  `json:"index_name"`
}

type DrbtPrcRnk struct {
	Wgt        float64 `json:"weight"`
	TmStmp     int64   `json:"timestamp"`
	Prc        float64 `json:"price"`
	Identifier string  `json:"identifier"`
	Enabled    bool    `json:"enabled"`
}

type NtfctnDrbtPrcRnk []DrbtPrcRnk

type NtfctnEstmtdExprtnPrc struct {
	Seconds     int     `json:"seconds"`
	Prc         float64 `json:"price"`
	IsEstimated bool    `json:"is_estimated"`
}

type MrkPrcOptn struct {
	IV         float64 `json:"iv"`
	Instrument string  `json:"instrument_name"`
	Mark       float64 `json:"mark_price"`
}

type NtfctnMrkPrcOptns []MrkPrcOptn

type NtfctnOrdrBkGrp struct {
	TmStmp     int64       `json:"timestamp"`
	Instrument string      `json:"instrument_name"`
	ChgID      int64       `json:"change_id"`
	Bids       [][]float64 `json:"bids"` // [price, amount]
	Asks       [][]float64 `json:"asks"` // [price, amount]
}

// NtfctnItmOrdrBk ...
// ["change",6947.0,82640.0]
// ["new",6942.5,6940.0]
// ["delete",6914.0,0.0]
type NtfctnItmOrdrBk struct {
	Actn string  `json:"action"`
	Prc  float64 `json:"price"`
	Amt  float64 `json:"amount"`
}

func (item *NtfctnItmOrdrBk) UnmarshalJSON(b []byte) error {
	// b: ["new",59786.0,10.0]
	// log.Printf("b=%v", string(b))
	s := strings.TrimRight(strings.TrimLeft(string(b), "["), "]")
	l := strings.Split(s, ",")
	if len(l) != 3 {
		return fmt.Errorf("fail to UnmarshalJSON [%v]", string(b))
	}
	item.Actn = strings.ReplaceAll(l[0], `"`, "")
	item.Prc, _ = strconv.ParseFloat(l[1], 64)
	item.Amt, _ = strconv.ParseFloat(l[2], 64)
	return nil
}

type NtfctnOrdrBk struct {
	Type       string            `json:"type"`
	TmStmp     int64             `json:"timestamp"`
	Instrument string            `json:"instrument_name"`
	PrvChgID   int64             `json:"prev_change_id"`
	ChgID      int64             `json:"change_id"`
	Bids       []NtfctnItmOrdrBk `json:"bids"` // [action, price, amount]
	Asks       []NtfctnItmOrdrBk `json:"asks"` // [action, price, amount]
}

type NtfctnOrdrBkRaw struct {
	TmStmp     int64             `json:"timestamp"`
	Instrument string            `json:"instrument_name"`
	PrvChgID   int64             `json:"prev_change_id"`
	ChgID      int64             `json:"change_id"`
	Bids       []NtfctnItmOrdrBk `json:"bids"` // [action, price, amount]
	Asks       []NtfctnItmOrdrBk `json:"asks"` // [action, price, amount]
}

type NtfctnPrptl struct {
	Interest float64 `json:"interest"`
}

type NtfctnPrtflio struct {
	TotalPl           float64 `json:"total_pl"`
	SessionUpl        float64 `json:"session_upl"`
	SessionRpl        float64 `json:"session_rpl"`
	SessionFunding    float64 `json:"session_funding"`
	PrtflioMrgnEnbld  bool    `json:"portfolio_margining_enabled"`
	OptnsVega         float64 `json:"options_vega"`
	OptnsTheta        float64 `json:"options_theta"`
	OptnsSessionUpl   float64 `json:"options_session_upl"`
	OptnsSessionRpl   float64 `json:"options_session_rpl"`
	OptnsPl           float64 `json:"options_pl"`
	OptnsGamma        float64 `json:"options_gamma"`
	OptnsDelta        float64 `json:"options_delta"`
	MrgnBal           float64 `json:"margin_balance"`
	MaintMrgn         float64 `json:"maintenance_margin"`
	InitMrgn          float64 `json:"initial_margin"`
	FuturesSessionUpl float64 `json:"futures_session_upl"`
	FuturesSessionRpl float64 `json:"futures_session_rpl"`
	FuturesPl         float64 `json:"futures_pl"`
	Equity            float64 `json:"equity"`
	DeltaTotal        float64 `json:"delta_total"`
	Ccy               string  `json:"currency"`
	Bal               float64 `json:"balance"`
	AvailableWdrFunds float64 `json:"available_withdrawal_funds"`
	AvailableFunds    float64 `json:"available_funds"`
}

type NtfctnQut struct {
	TmStmp     int64   `json:"timestamp"`
	Instrument string  `json:"instrument_name"`
	BestBid    float64 `json:"best_bid_price"`
	BestBidAmt float64 `json:"best_bid_amount"`
	BestAsk    float64 `json:"best_ask_price"`
	BestAskAmt float64 `json:"best_ask_amount"`
}

type NtfctnTckr struct {
	TmStmp int64 `json:"timestamp"`
	Stats  struct {
		Volume float64 `json:"volume"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
	} `json:"stats"`
	State          string  `json:"state"`
	SttlmntPrc     float64 `json:"settlement_price"`
	OpnIntrst      float64 `json:"open_interest"`
	MinPrc         float64 `json:"min_price"`
	MaxPrc         float64 `json:"max_price"`
	Mark           float64 `json:"mark_price"`
	Last           float64 `json:"last_price"`
	Instrument     string  `json:"instrument_name"`
	IndxPrc        float64 `json:"index_price"`
	Funding8H      float64 `json:"funding_8h"`
	CurrentFunding float64 `json:"current_funding"`
	BestBid        float64 `json:"best_bid_price"`
	BestBidAmt     float64 `json:"best_bid_amount"`
	BestAsk        float64 `json:"best_ask_price"`
	BestAskAmt     float64 `json:"best_ask_amount"`
}

type NtfctnTrades []Trade

type NtfctnUserTrades []UserTrade

type NtfctnUserOrdr []Order

type NtfctnUserChgs struct {
	Trades []UserTrade `json:"trades"`
	Pstns  []PosnOut   `json:"positions"`
	Orders []Order     `json:"orders"`
}
