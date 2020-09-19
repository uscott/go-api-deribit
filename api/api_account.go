package api

import (
	"github.com/uscott/go-api-deribit/inout"
)

// GetAccountSummary requests "private/get_account_summary"
func (c *Client) GetAccountSummary(ccy string, extended bool, result *inout.AcctSummaryOut) error {
	return c.Call(
		"private/get_account_summary",
		&inout.AcctSummaryIn{Ccy: ccy, Extended: extended},
		result,
	)
}

// GetPositionInstrument requests "private/get_position"
func (c *Client) GetPositionInstrument(contract string, result *inout.PosnOut) error {
	return c.Call(
		"private/get_position",
		&inout.PosnInstrmtIn{Instrument: contract},
		result)
}

// GetPositionCurrency requests "private/get_positions"
func (c *Client) GetPositionCurrency(
	params *inout.PosnCcyIn) (result []inout.PosnOut, err error) {

	err = c.Call("private/get_positions", params, &result)
	return
}
