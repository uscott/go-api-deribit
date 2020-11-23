package api

import (
	"github.com/uscott/go-api-deribit/inout"
)

// GetHistoricalVolatility returns slice of []float64 where each []float64 is of length 2
// with the first element being the exchange time stamp and the 2nd the historical vol in
// percent points as measured at the time stamp
func (c *Client) GetHistoricalVolatility(currency string) (vols [][]float64, err error) {
	arg := inout.HistVolIn{Ccy: currency}
	err = c.Call("public/get_historical_volatility", &arg, &vols)
	return
}
