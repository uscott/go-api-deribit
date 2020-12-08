package api

import (
	"fmt"
	"math"
	"time"

	"github.com/uscott/go-api-deribit/inout"
	"github.com/uscott/go-tools/errs"
)

var contractMap = map[string]string{
	"JAN": "F",
	"FEB": "G",
	"MAR": "H",
	"APR": "J",
	"MAY": "K",
	"JUN": "M",
	"JUL": "N",
	"AUG": "Q",
	"SEP": "U",
	"OCT": "V",
	"NOV": "X",
	"DEC": "Z",
}

const exchTmStmpUnit = time.Millisecond

// ConvertExchStmp converts an exchange time stamp
// to a client-side time.Time
func ConvertExchStmp(ts int64) time.Time {
	ts *= int64(exchTmStmpUnit) / int64(time.Nanosecond)
	return time.Unix(ts/int64(time.Second), ts%int64(time.Second)).UTC()
}

// ExchangeTime returns the exchange time as
// a client-side time.Time
func (c *Client) ExchangeTime() (time.Time, error) {
	var (
		ms  int64
		err error
	)
	if ms, err = c.GetTime(); err != nil {
		return time.Time{}, err
	}
	return ConvertExchStmp(ms), nil
}

// GetTime retrieves the exchange server time
func (c *Client) GetTime() (result int64, err error) {
	err = c.Call("public/get_time", nil, &result)
	return
}

// Test call "public/test"
func (c *Client) Test() (result inout.TestOut, err error) {
	err = c.Call("public/test", nil, &result)
	return
}

// ConvertSymbol converts the exchange contract symbol
// to one resembling a more traditional format
func ConvertSymbol(s string) string {
	var phys, day, month, year string
	switch {
	case s == "BTC":
		return s
	case len(s) > 4 && s[4:] == "PERPETUAL":
		return s
	case len(s) > 15 && s[0:5] == "BASIS":
		phys = s[0:10]
		day = s[10:12]
		month = s[12:15]
		year = s[15:]
	case len(s) > 9:
		phys = s[0:4]
		day = s[4:6]
		month = s[6:9]
		year = s[9:]
	default:
		return ""
	}
	// return phys + "20" + year + contractMap[month] + day
	return fmt.Sprintf("%v20%v-%v-%v", phys, year, contractMap[month], day)
}

// Inverse returns 1/x expect when abs(x) is <= Small in which case it
// returns +/- Big
func Inverse(x float64) float64 {
	switch {
	case math.Abs(x) > SMALL:
		return 1.0 / x
	case x > 0:
		return BIG
	default:
		return -BIG
	}
}

func CpyBook(src *inout.BookOut, dst *Book) error {
	if src == nil {
		return fmt.Errorf("source pointer is nil")
	}
	if dst == nil {
		return fmt.Errorf("destination pointer is nil")
	}
	dst.Contract = src.Instrument
	dst.IndxPrc = src.IndxPrc
	dst.Last = src.Last
	dst.OpenInterest = src.OpnIntrst
	dst.Stats = src.Stats
	dst.TimeStamp = ConvertExchStmp(src.TmStmp)
	nbids, nasks := len(src.Bids), len(src.Asks)
	if cap(dst.Bids) < nbids {
		dst.Bids = make([]Quote, nbids)
	}
	dst.Bids = dst.Bids[:nbids]
	if cap(dst.Asks) < nasks {
		dst.Asks = make([]Quote, nasks)
	}
	dst.Asks = dst.Asks[:nasks]
	q := dst.Bids
	for i, bid := range src.Bids {
		q[i].Prc, q[i].Amt = bid[0], bid[1]
	}
	q = dst.Asks
	for i, ask := range src.Asks {
		q[i].Prc, q[i].Amt = ask[0], ask[1]
	}
	if nbids > 0 {
		dst.BestBid = dst.Bids[0]
	} else {
		dst.BestBid = Quote{Prc: math.NaN(), Amt: math.NaN()}
	}
	if nasks > 0 {
		dst.BestAsk = dst.Asks[0]
	} else {
		dst.BestAsk = Quote{Prc: math.NaN(), Amt: math.NaN()}
	}
	return nil
}

func (c *Client) GetOneInstrument(
	contract, currency, kind string) (inout.InstrumentOut, error) {

	instruments, err := c.GetInstruments(currency, kind, false)
	if err != nil {
		return inout.InstrumentOut{}, err
	}
	for _, ins := range instruments {
		if ins.Instrument == contract {
			return ins, nil
		}
	}
	return inout.InstrumentOut{}, fmt.Errorf("could not get instrument")
}

func NewBook(bk *inout.BookOut, expiration time.Time) (*Book, error) {
	if bk == nil {
		return nil, errs.ErrNilPtr
	}
	out := Book{Expiration: expiration}
	err := CpyBook(bk, &out)
	return &out, err
}

func (c *Client) NewBook(contract, currency, kind string, depth int) (*Book, error) {
	con, ccy := contract, currency
	ins, err := c.GetOneInstrument(con, ccy, "future")
	if err != nil {
		return nil, err
	}
	var bkraw inout.BookOut
	err = c.GetBook(con, depth, &bkraw)
	if err != nil {
		return nil, err
	}
	expi := ConvertExchStmp(ins.ExprtnTmStmp)
	return NewBook(&bkraw, expi)
}

// NewFuturesData returns an allocated pointer to a FuturesData struct
// based on the data in the InstrumentOut argument
func NewFuturesData(ins *inout.InstrumentOut) (*FuturesData, error) {
	if ins == nil {
		return nil, errs.ErrNilPtr
	}
	f := FuturesData{InstrumentOut: *ins, Expiration: time.Time{}, IsSwap: false}
	f.IsSwap = ins.StlmntPrd == "perpetual"
	if f.IsSwap {
		f.Expiration = time.Date(9999, 12, 31, 0, 0, 0, 0, time.FixedZone("utc", 0))
	} else {
		f.Expiration = ConvertExchStmp(ins.ExprtnTmStmp)
	}
	return &f, nil
}

func (c *Client) NewFuturesData(contract, currency string) (*FuturesData, error) {
	ins, err := c.GetOneInstrument(contract, currency, "future")
	if err != nil {
		return nil, err
	}
	return NewFuturesData(&ins)
}

func PruneOrdersFromBook(bk *Book, orders []inout.Order) error {
	if bk == nil {
		return errs.ErrNilPtr
	}
	bidOrders := make([]inout.Order, 0, len(orders))
	askOrders := make([]inout.Order, 0, len(orders))
	for _, o := range orders {
		switch o.Direction {
		case DirBuy:
			if o.Instrument == bk.Contract {
				bidOrders = append(bidOrders, o)
			}
		case DirSell:
			if o.Instrument == bk.Contract {
				askOrders = append(askOrders, o)
			}
		default:
			return fmt.Errorf("unknown order direction: %+v", o)
		}
	}
	for i := range bk.Bids {
		err := PruneOrdersFromQuote(&bk.Bids[i], &bidOrders)
		if err != nil {
			return err
		}
	}
	for i := range bk.Asks {
		err := PruneOrdersFromQuote(&bk.Asks[i], &askOrders)
		if err != nil {
			return err
		}
	}
	RmZeroQuotes(&bk.Bids)
	RmZeroQuotes(&bk.Asks)
	if len(bk.Bids) > 0 {
		bk.BestBid = bk.Bids[0]
	} else {
		bk.BestBid = Quote{Prc: math.NaN(), Amt: math.NaN()}
	}
	if len(bk.Asks) > 0 {
		bk.BestAsk = bk.Asks[0]
	} else {
		bk.BestAsk = Quote{Prc: math.NaN(), Amt: math.NaN()}
	}
	return nil
}

// PruneOrdersFromQuote0 is used to remove user orders from the order book
func PruneOrdersFromQuote0(quote *Quote, orders *[]inout.Order, maxdiff float64) error {
	if quote == nil {
		return errs.ErrNilPtr
	}
	if len(*orders) == 0 {
		return nil
	}
	rmvIndcs := make([]int, 0, len(*orders))
	prc, amt, mchAmt := quote.Prc, quote.Amt, 0.0
	for i, o := range *orders {
		if math.Abs(prc-float64(o.Prc)) < maxdiff {
			mchAmt += o.Amt
			rmvIndcs = append(rmvIndcs, i)
		}
	}
	nRmvd := 0
	for _, i := range rmvIndcs {
		j := i - nRmvd
		switch j {
		case 0:
			*orders = (*orders)[1:]
		case len(*orders) - 1:
			*orders = (*orders)[0:j]
		default:
			*orders = append((*orders)[:j], (*orders)[j+1:]...)
		}
		nRmvd++
	}
	quote.Amt = amt - mchAmt
	return nil
}

func PruneOrdersFromQuote(quote *Quote, orders *[]inout.Order) error {
	if quote == nil {
		return errs.ErrNilPtr
	}
	if len(*orders) == 0 {
		return nil
	}
	rmv := make([]int, 0, len(*orders))
	prc, amt, mchAmt := quote.Prc, quote.Amt, 0.0
	for i, o := range *orders {
		if math.Abs(prc-float64(o.Prc)) < SMALL {
			mchAmt += o.Amt
			rmv = append(rmv, i)
		}
	}
	nrm := 0
	for _, i := range rmv {
		rmOrder(orders, i-nrm)
		nrm++
	}
	quote.Amt = amt - mchAmt
	return nil
}

func rmOrder(orders *[]inout.Order, i int) {
	if orders == nil {
		return
	}
	n := len(*orders)
	if i < 0 || n <= i || n == 0 {
		return
	}
	switch i {
	case 0:
		*orders = (*orders)[1:]
	case n - 1:
		*orders = (*orders)[:n-1]
	default:
		*orders = append((*orders)[:i], (*orders)[i+1:]...)
	}
}

func rmQuote(quotes *[]Quote, i int) {
	if quotes == nil {
		return
	}
	n := len(*quotes)
	if i < 0 || n <= i || n == 0 {
		return
	}
	switch i {
	case 0:
		*quotes = (*quotes)[1:]
	case n - 1:
		*quotes = (*quotes)[:n-1]
	default:
		*quotes = append((*quotes)[:i], (*quotes)[i+1:]...)
	}
}

func RmZeroQuotes(quotes *[]Quote) {
	if quotes == nil {
		return
	}
	n := len(*quotes)
	if n == 0 {
		return
	}
	rmv := make([]int, 0, n)
	for i, q := range *quotes {
		if q.Amt < math.SmallestNonzeroFloat64 {
			rmv = append(rmv, i)
		}
	}
	nrm := 0
	for _, i := range rmv {
		rmQuote(quotes, i-nrm)
		nrm++
	}
}
