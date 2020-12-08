package api

import (
	"math"
	"time"

	"github.com/uscott/go-api-deribit/inout"
	"github.com/uscott/go-tools/tmath"
)

var (
	imax  = tmath.Imax
	isinf = math.IsInf
	isnan = math.IsNaN
)

// BookDepth specifies the book depth for queries
const BookDepth = 50

type rqstCntData struct {
	mch int
	non int
}

type rqstTmrData struct {
	t0 time.Time
	t1 time.Time
	dt time.Duration
}

// AccountBalance contains the current and initial
// mark to market and theoretical account balances
// in both USD and BTC
type AccountBalance struct {
	Crnt struct {
		Mrkt Balance
		Theo Balance
	}
	Init struct {
		Mrkt Balance
		Theo Balance
	}
}

// Balance contains balance data for USD and BTC
type Balance struct {
	USD float64
	Ccy float64
}

// FuturesData embeds struct inout.InstrumentOut
// as well as time to expiration and IsSwap
type FuturesData struct {
	inout.InstrumentOut
	Expiration time.Time
	IsSwap     bool
}

// Quote has fields prc, amt and direction
// corresonding to an order book quote
type Quote struct {
	Prc float64 `json:"price"`
	Amt float64 `json:"amount"`
}

// Book contains orderbook bids, asks with
// user orders pruned out
type Book struct {
	BestBid      Quote           `json:"best_bid"`
	BestAsk      Quote           `json:"best_ask"`
	Bids         []Quote         `json:"bids"`
	Asks         []Quote         `json:"asks"`
	Contract     string          `json:"contract"`
	Expiration   time.Time       `json:"expiration"`
	IndxPrc      float64         `json:"index_price"`
	Last         float64         `json:"last_price"`
	OpenInterest float64         `json:"open_interest"`
	Stats        inout.TckrStats `json:"stats"`
	TimeStamp    time.Time       `json:"timestamp"`
}

// InitBook sets Bids and Asks fields
func InitBook(ob *Book) {
	ob.Bids = make([]Quote, BookDepth)
	ob.Asks = make([]Quote, BookDepth)
}
