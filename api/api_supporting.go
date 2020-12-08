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

// NewFuturesData returns an allocated pointer to a FuturesData struct
// based on the data in the InstrumentOut argument
func NewFuturesData(ins *inout.InstrumentOut) (*FuturesData, error) {
	if ins == nil {
		return nil, errs.ErrNilPtr
	}
	f := FuturesData{InstrumentOut: *ins, Expiration: time.Time{}, IsSwap: false}
	f.Expiration = ConvertExchStmp(ins.ExprtnTmStmp)
	f.IsSwap = ins.StlmntPrd == "perpetual"
	return &f, nil
}

// PruneUsrOrdrsFromQuts is used to remove user orders from the order book
func PruneUsrOrdrsFromQuts(quotes []float64, orders *[]inout.Order, maxdiff float64) {
	if len(*orders) == 0 {
		return
	}
	rmvIndcs := make([]int, 0, len(*orders))
	prc, amt, mchAmt := quotes[0], quotes[1], 0.0
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
	quotes[1] = amt - mchAmt
}
