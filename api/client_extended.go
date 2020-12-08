package api

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/uscott/go-api-deribit/inout"
	"github.com/uscott/go-tools/errs"
)

// ClientExtended extendes the base Client
type ClientExtended struct {
	*Client
	Blnc      AccountBalance
	BkSummary map[string]inout.BkSummaryOut
	Buffer    *bytes.Buffer
	ConSz     map[string]float64
	Contracts []string
	Deltas    map[string]float64
	Filepath  struct {
		Base   string
		Dir    string
		TmStmp time.Time
		Path   string
	}
	Futures     map[string]FuturesData
	Instruments []inout.InstrumentOut
	IsSwap      map[string]bool
	ObAdj       map[string]Book
	ObRaw       map[string]inout.BookOut
	OpnOrdrs    map[string][]inout.Order
	PctDelta    float64
	Posn        map[string]inout.PosnOut
	Spot        float64
	SymbolMap   struct {
		Fwd map[string]string
		Inv map[string]string
	}
	Symbols    []string
	Tckr       map[string]inout.TckrOut
	TckSz      map[string]float64
	TmClnt     time.Time
	TmExch     time.Time
	UserTrades map[string][]inout.UserTrade
}

// NewClientMin makes a pointer to a new ClientExtended
// with a minimum of initiatialization
func NewClientMin(cfg *Configuration) (*ClientExtended, error) {
	var (
		c   ClientExtended
		err error
	)
	if c.Client, err = New(cfg); err != nil {
		return &c, err
	}
	c.Buffer = bytes.NewBuffer(make([]byte, 512))
	c.Futures = make(map[string]FuturesData)
	c.Deltas = make(map[string]float64)
	c.BkSummary = make(map[string]inout.BkSummaryOut)
	c.Posn = make(map[string]inout.PosnOut)
	c.OpnOrdrs = make(map[string][]inout.Order)
	c.ObRaw = make(map[string]inout.BookOut)
	c.ObAdj = make(map[string]Book)
	c.IsSwap = make(map[string]bool)
	c.TckSz = make(map[string]float64)
	c.Tckr = make(map[string]inout.TckrOut)
	c.ConSz = make(map[string]float64)
	c.UserTrades = make(map[string][]inout.UserTrade)
	return &c, nil
}

// NewClientExtended returns pointer to a new ClientExtended
func NewClientExtended(cfg *Configuration) (*ClientExtended, error) {
	c, err := NewClientMin(cfg)
	if err = c.UpdtFutures(); err != nil {
		return c, err
	}
	for k, v := range c.Futures {
		c.IsSwap[k] = v.IsSwap
		c.TckSz[k] = v.TckSz
		c.ConSz[k], err = c.GetContractSize(k)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func (c *ClientExtended) addSubs(
	channels *[]string,
	prefix string,
	suffix string) {
	for _, s := range c.Contracts {
		c.Buffer.Reset()
		if len(prefix) > 0 {
			c.Buffer.WriteString(prefix)
			c.Buffer.WriteString(".")
		}
		c.Buffer.WriteString(s)
		if len(suffix) > 0 {
			c.Buffer.WriteString(".")
			c.Buffer.WriteString(suffix)
		}
		*channels = append(*channels, c.Buffer.String())
	}
	c.Buffer.Reset()
}

// CreateSymbols creates the Symbols, Contracts and SymbolMap fields
func (c *ClientExtended) CreateSymbols() {
	c.Contracts = make([]string, len(c.Futures))
	i := 0
	for k := range c.Futures {
		c.Contracts[i] = k
		i++
	}
	sort.StringSlice(c.Contracts).Sort()
	c.Symbols = make([]string, len(c.Contracts)+1)
	c.Symbols[0] = c.Config.Currency
	copy(c.Symbols[1:], c.Contracts)
	c.SymbolMap.Fwd = make(map[string]string)
	c.SymbolMap.Inv = make(map[string]string)
	for _, k := range c.Symbols {
		h := ConvertSymbol(k)
		c.SymbolMap.Fwd[k] = h
		c.SymbolMap.Inv[h] = k
	}
}

// MakeSubscriptions makes the slice of channels to subscribe to
func (c *ClientExtended) MakeSubscriptions() []string {
	channels := make([]string, 0, 128)
	// Book
	// book.(instrument_name).(interval/"raw")
	c.addSubs(&channels, "book", "raw")
	// Price index
	// deribit_price_index.(index_name)
	channels = append(channels, "deribit_price_index.btc_usd")
	// Platform state
	channels = append(channels, "platform_state")
	// Quotes
	// quote.(instrument name)
	c.addSubs(&channels, "quote", "")
	// Tickers:
	// ticker.(intrument_name).(interval/"raw")
	c.addSubs(&channels, "ticker", "raw")
	// User orders:
	// user.orders.(instrument_name).(interval/"raw")
	c.addSubs(&channels, "user.orders", "raw")
	// User portfolio:
	// user.portfolio.(currency)
	channels = append(channels, fmt.Sprintf("user.portfolio.%v", c.Config.Currency))
	// User trades:
	// user.trades.(instrument_name).(interval/"raw")
	c.addSubs(&channels, "user.trades", "raw")
	return channels
}

// UpdtAcct updates account summary
func (c *ClientExtended) UpdtAcct() error {
	return c.GetAccountSummary(c.Config.Currency, true, &c.Acct)
}

// UpdtBalance updates balance in BTC
// and USD terms
func (c *ClientExtended) UpdtBalance() error {
	e, s := c.Acct.Equity, c.Spot
	switch {
	case isnan(e) || isnan(s):
		return errs.ErrNaN
	case isinf(e, 0) || isinf(s, 0):
		return errs.ErrInf
	case s == 0:
		return errs.ErrDivByZero
	default:
	}
	c.Blnc.Crnt.Mrkt.USD = e * s
	c.Blnc.Crnt.Mrkt.Ccy = e
	c.Blnc.Crnt.Theo.USD = c.Blnc.Crnt.Mrkt.USD
	for _, v := range c.Posn {
		if v.SzCcy != 0.0 {
			diff := v.IndxPrc - v.Mark
			c.Blnc.Crnt.Theo.USD += v.SzCcy * diff
		}
	}
	c.Blnc.Crnt.Theo.Ccy = c.Blnc.Crnt.Theo.USD / s
	c.Deltas[c.Config.Currency] = c.Blnc.Crnt.Theo.Ccy
	return nil
}

// UpdtBkSummary updates book summary for contract
func (c *ClientExtended) UpdtBkSummary(contract string) error {
	bkSummary, err := c.GetBookSummaryByInstrument(
		&inout.BkSummaryByInstrmtIn{Instrument: contract})
	if err != nil {
		return err
	}
	c.BkSummary[contract] = bkSummary[0]
	return nil
}

// UpdtFutures updates futures info map
func (c *ClientExtended) UpdtFutures() (err error) {
	ccy := c.Config.Currency
	c.Instruments, err = c.GetInstruments(
		&inout.InstrumentIn{Ccy: ccy, Kind: "future", Expired: false})
	if err != nil {
		return err
	}
	for _, i := range c.Instruments {
		if i.Kind == "future" {
			c.Futures[i.Instrument] = FuturesData{i, time.Time{}, false}
		}
	}
	c.CreateSymbols()
	for k, v := range c.Futures {
		v.Expiration = ConvertExchStmp(v.ExprtnTmStmp)
		v.IsSwap = v.StlmntPrd == "perpetual"
		c.Futures[k] = v
	}
	return nil
}

// UpdtOpnOrdrs updates the Client's field of open orders
func (c *ClientExtended) UpdtOpnOrdrs(contract string) error {
	orders, err := c.GetOpenOrdersByInstrument(
		&inout.OpnOrdrsByInstrmtIn{Instrument: contract})
	if err != nil {
		return err
	}
	c.OpnOrdrs[contract] = orders
	return nil
}

// UpdtOrdrBkAdj updates the order book and prunes user orders
func (c *ClientExtended) UpdtOrdrBkAdj(contract string) error {
	k := contract
	maxDiff := c.TckSz[k] / 4.0
	bids := make([][]float64, len(c.ObRaw[k].Bids))
	asks := make([][]float64, len(c.ObRaw[k].Asks))
	for i, qut := range c.ObRaw[k].Bids {
		bids[i] = []float64{qut[0], qut[1]}
	}
	for i, qut := range c.ObRaw[k].Asks {
		asks[i] = []float64{qut[0], qut[1]}
	}
	bidOrdrs := make([]inout.Order, 0, len(c.OpnOrdrs[k]))
	askOrdrs := make([]inout.Order, 0, len(c.OpnOrdrs[k]))
	for _, o := range c.OpnOrdrs[k] {
		if o.Drctn == DrctnBuy {
			bidOrdrs = append(bidOrdrs, o)
		} else if o.Drctn == DrctnSell {
			askOrdrs = append(askOrdrs, o)
		} else {
			return fmt.Errorf("order direction is neither buy nor sell: %v", o)
		}
	}
	for _, qut := range bids {
		PruneUsrOrdrsFromQuts(qut, &bidOrdrs, maxDiff)
	}
	var bestBidAmt, bestBid float64
	bidsAdj := make([]Quote, 0, len(bids))
	for _, qut := range bids {
		prc, amt := qut[0], qut[1]
		if amt > 0 {
			if prc > bestBid {
				bestBid, bestBidAmt = prc, amt
			}
			bidsAdj = append(bidsAdj, Quote{Amt: amt, Prc: prc})
		}
	}
	for _, qut := range asks {
		PruneUsrOrdrsFromQuts(qut, &askOrdrs, maxDiff)
	}
	var bestAskAmt float64
	bestAsk, asksAdj := math.MaxFloat64, make([]Quote, 0, len(asks))
	for _, qut := range asks {
		prc, amt := qut[0], qut[1]
		if amt > 0 {
			if prc < bestAsk {
				bestAsk, bestAskAmt = prc, amt
			}
			asksAdj = append(asksAdj, Quote{Amt: amt, Prc: prc})
		}
	}
	c.ObAdj[k] = Book{
		TimeStamp: ConvertExchStmp(c.ObRaw[k].TmStmp),
		BestBid:   Quote{Prc: bestBid, Amt: bestBidAmt},
		BestAsk:   Quote{Prc: bestAsk, Amt: bestAskAmt},
		Bids:      bidsAdj,
		Asks:      asksAdj,
	}
	return nil
}

// UpdtBookRaw updates client's copy of order book
func (c *ClientExtended) UpdtBookRaw(contract string) (e error) {
	var ob inout.BookOut
	if e = c.GetBook(contract, BookDepth, &ob); e != nil {
		return e
	}
	c.ObRaw[contract] = ob
	return nil
}

// UpdtPctDelta computes the % delta and
// updates corresponding field
func (c *ClientExtended) UpdtPctDelta() error {
	d := 0.0
	for _, v := range c.Deltas {
		d += v
	}
	b := c.Blnc.Crnt.Theo.Ccy
	switch {
	case isnan(b) || isnan(d):
		return errs.ErrNaN
	case isinf(b, 0) || isinf(d, 0):
		return errs.ErrInf
	case b == 0.0:
		return errs.ErrDivByZero
	default:
		c.PctDelta = d / b
	}
	return nil
}

// UpdtPosition updates the position info for the contract
func (c *ClientExtended) UpdtPosition(contract string) (e error) {
	posn := inout.PosnOut{}
	e = c.GetPositionInstrument(contract, &posn)
	if e != nil {
		return e
	}
	c.Posn[contract] = posn
	if c.Config.DebugMode {
		fmt.Printf("%v\n", posn.Delta)
	}
	c.Deltas[contract] = posn.Delta
	return nil
}

// UpdtSpot updates spot price and corresponding field
func (c *ClientExtended) UpdtSpot() (e error) {
	c.Spot, e = c.GetIndex(c.Config.Currency)
	if e != nil {
		return e
	}
	return nil
}

// UpdtStatus updates most of the client fields
func (c *ClientExtended) UpdtStatus() (e error) {
	if e = c.UpdtAcct(); e != nil {
		return e
	}
	if e = c.UpdtSpot(); e != nil {
		return e
	}
	c.TmClnt = time.Now().UTC()
	c.TmExch, e = c.ExchangeTime()
	if e != nil {
		return e
	}
	for _, k := range c.Contracts {
		if e = c.UpdtBkSummary(k); e != nil {
			return e
		}
		if e = c.UpdtPosition(k); e != nil {
			return e
		}
	}
	if e = c.UpdtBalance(); e != nil {
		return e
	}
	return nil
}

// UpdtTckr updates ticker info for the contract
func (c *ClientExtended) UpdtTckr(contract string) (e error) {
	var tckr inout.TckrOut
	if e = c.GetTckr(contract, &tckr); e != nil {
		return e
	}
	c.Tckr[contract] = tckr
	return nil
}

// UpdtUserTrades updates user trades
func (c *ClientExtended) UpdtUserTrades(
	contract string, count int) error {

	params := inout.TradesByInstrmtIn{
		Instrument: contract,
		Count:      count,
		IncludeOld: true,
	}
	out := inout.UserTradesOut{}
	err := c.GetUserTradesByInstrument(&params, &out)
	if err == nil {
		c.UserTrades[contract] = out.Trades
		return nil
	}
	return err
}
