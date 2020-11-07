package api

import (
	"time"

	"github.com/uscott/go-api-deribit/inout"
	"github.com/uscott/go-tools/errs"
)

var timeZero = time.Date(1970, 1, 1, 0, 0, 0, 0, time.FixedZone("utc", 0))

// Buy posts a buy order to the exchange
func (c *Client) Buy(params *inout.OrderIn, result *inout.OrderOut) error {
	return c.Call("private/buy", params, result)
}

// Sell posts a sell order
func (c *Client) Sell(params *inout.OrderIn, result *inout.OrderOut) error {
	return c.Call("private/sell", params, result)
}

// Edit posts an edit request
func (c *Client) Edit(params *inout.EditIn, result *inout.OrderOut) error {
	return c.Call("private/edit", params, result)
}

// Trade posts a buy or sell depending on the sign of amount
// Buy for positive, sell for negative
func (c *Client) Trade(params *inout.OrderIn, result *inout.OrderOut) (e error) {
	if params == nil {
		return errs.ErrNilPtr
	}
	if params.Amt < 0 {
		params.Amt = -params.Amt
		e = c.Call("private/sell", params, result)
		params.Amt = -params.Amt
	} else {
		e = c.Call("private/buy", params, result)
	}
	return
}

// Cancel submits a cancel request
func (c *Client) Cancel(oid string, result *inout.Order) error {
	return c.Call(
		"private/cancel",
		&inout.CancelOrderIn{OrderID: oid},
		result)
}

// CancelAll requests to cancel all open orders
func (c *Client) CancelAll() (result int, err error) {
	err = c.Call("private/cancel_all", nil, &result)
	return
}

// CancelAllByCurrency requests to cancell all orders for a given currency
func (c *Client) CancelAllByCurrency(
	params *inout.CancelAllByCcyIn) (result int, err error) {
	err = c.Call("private/cancel_all_by_currency", params, &result)
	return
}

// CancelAllByInstrument requests to cancel all orders for a given
// instrument/contract
func (c *Client) CancelAllByInstrument(
	params *inout.CancelAllByInstrmtIn) (result int, err error) {
	err = c.Call("private/cancel_all_by_instrument", params, &result)
	return
}

func (c *Client) GetOpenOrdersByCurrency(
	params *inout.OpnOrdrsByCcyIn) (result []inout.Order, err error) {
	err = c.Call("private/get_open_orders_by_currency", params, &result)
	return
}

func (c *Client) GetOpenOrdersByInstrument(
	params *inout.OpnOrdrsByInstrmtIn) (result []inout.Order, err error) {
	err = c.Call("private/get_open_orders_by_instrument", params, &result)
	return
}

func (c *Client) GetStopOrderHistory(
	params *inout.StopOrderHistoryIn, result *inout.StopOrderHistoryOut) (err error) {
	err = c.Call("private/get_stop_order_history", params, result)
	return
}

func (c *Client) GetUserTradesByCurrency(
	params *inout.TradesByCcyIn, result *inout.UserTradesOut) error {

	return c.Call("private/get_user_trades_by_currency", params, result)
}

func (c *Client) GetUserTradesByCurrencyAndTime(
	params *inout.TradesByCcyAndTmIn, result *inout.UserTradesOut) error {

	return c.Call("private/get_user_trades_by_currency_and_time", params, result)
}

func (c *Client) GetUserTradesByInstrument(
	params *inout.TradesByInstrmtIn, result *inout.UserTradesOut) error {

	return c.Call("private/get_user_trades_by_instrument", params, result)

}

func (c *Client) GetUserTradesByInstrumentAndTime(
	params *inout.TradesByInstrmtAndTmIn, result *inout.UserTradesOut) error {

	return c.Call("private/get_user_trades_by_instrument_and_time", params, result)
}

func (c *Client) GetUserTradesByInstrumentAndTimeExtended(
	instrument string, start, end time.Time) ([]inout.UserTrade, error) {

	startStamp := int64(start.Sub(timeZero) / time.Millisecond)
	endStamp := int64(end.Sub(timeZero) / time.Millisecond)
	out := inout.UserTradesOut{}
	var err error
	const step = 250
	for cnt := step; cnt <= 1000; cnt += step {
		err = c.GetUserTradesByInstrumentAndTime(
			&inout.TradesByInstrmtAndTmIn{
				Instrument:  instrument,
				StartTmStmp: startStamp,
				EndTmStmp:   endStamp,
				Count:       cnt,
				IncludeOld:  true,
			},
			&out)
		if err != nil {
			return []inout.UserTrade{}, err
		}
		if !out.HasMore {
			break
		}
		time.Sleep(time.Second)
	}
	return out.Trades, nil
}
