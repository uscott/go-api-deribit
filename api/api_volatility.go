package api

import (
	"github.com/uscott/go-api-deribit/inout"
)

func (c *Client) GetHistoricalVolatility(currency string) (vols interface{}, err error) {
	arg := inout.HistVolIn{Ccy: currency}
	err = c.Call("public/get_historical_volatility", &arg, &vols)
	return
}
