package api

import (
	"fmt"
	"math"

	"github.com/uscott/go-api-deribit/inout"
)

// GetBkSummaryByCurrency gets a book summary
func (c *Client) GetBookSummaryByCurrency(
	params *inout.BkSummaryByCcyIn) (result []inout.BkSummaryOut, err error) {
	err = c.Call("public/get_book_summary_by_currency", params, &result)
	return
}

// GetBkSummaryByInstrument gets book summary for given instrument
func (c *Client) GetBookSummaryByInstrument(
	params *inout.BkSummaryByInstrmtIn) (result []inout.BkSummaryOut, err error) {
	err = c.Call("public/get_book_summary_by_instrument", params, &result)
	return
}

// GetContractSize gets contract size for an instrument
func (c *Client) GetContractSize(contract string) (float64, error) {
	result := &inout.CntrctSzOut{}
	err := c.Call(
		"public/get_contract_size",
		&inout.CntrctSzIn{Instrument: contract},
		&result)
	return result.CntrctSz, err
}

// GetCurrencies gets spot currencies/indices on exchange
func (c *Client) GetCurrencies() (result []inout.Currency, err error) {
	err = c.Call("public/get_currencies", nil, &result)
	return
}

// GetIndex returns index value
func (c *Client) GetIndex(currency string) (float64, error) {
	result := inout.IndxOut{}
	err := c.Call("public/get_index", &inout.IndxIn{Ccy: currency}, &result)
	if err != nil {
		return math.NaN(), err
	}
	switch currency {
	case BTC:
		return result.BTC, nil
	case ETH:
		return result.ETH, nil
	case Edp:
		return result.EDP, nil
	default:
		return math.NaN(), fmt.Errorf("Unknown currency type")
	}
}

// GetInstruments returns instruments/contracts traded on exchange
// various related data
func (c *Client) GetInstruments(
	currency, kind string, expired bool) (result []inout.InstrumentOut, err error) {

	err = c.Call(
		"public/get_instruments",
		&inout.InstrumentIn{
			Ccy:     currency,
			Kind:    kind,
			Expired: expired,
		},
		&result)
	return result, err
}

func (c *Client) GetLastTradesByCurrency(
	params *inout.TradesByCcyIn) (result *inout.LastTradesOut, err error) {

	result = &inout.LastTradesOut{}
	err = c.Call("public/get_last_trades_by_currency", params, &result)
	return
}

func (c *Client) GetLastTradesByCurrencyAndTime(
	params *inout.TradesByCcyAndTmIn) (result *inout.LastTradesOut, err error) {

	result = &inout.LastTradesOut{}
	err = c.Call("public/get_last_trades_by_currency_and_time", params, &result)
	return
}

func (c *Client) GetLastTradesByInstrument(
	params *inout.TradesByInstrmtIn) (result *inout.LastTradesOut, err error) {

	result = &inout.LastTradesOut{}
	err = c.Call("public/get_last_trades_by_instrument", params, &result)
	return
}

func (c *Client) GetLastTradesByInstrumentAndTime(
	params *inout.TradesByInstrmtAndTmIn) (result *inout.LastTradesOut, err error) {

	result = &inout.LastTradesOut{}
	err = c.Call("public/get_last_trades_by_instrument_and_time", params, &result)
	return
}

// GetBook returns the order book for a given contract
func (c *Client) GetBook(
	contract string, depth int, result *inout.BookOut) error {

	return c.Call(
		"public/get_order_book",
		&inout.BookIn{Instrument: contract, Depth: depth},
		result)
}

// GetTckr gets the ticker data for a given contract
func (c *Client) GetTckr(contract string, result *inout.TckrOut) error {
	return c.Call(
		"public/ticker",
		&inout.TckrIn{Instrument: contract},
		result)
}
