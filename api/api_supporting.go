package api

import (
	"fmt"
	"math"

	"github.com/uscott/go-api-deribit/inout"
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
	case math.Abs(x) > Small:
		return 1.0 / x
	case x > 0:
		return Big
	default:
		return -Big
	}
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
