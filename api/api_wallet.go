package api

import (
	"github.com/uscott/go-api-deribit/inout"
)

// SubmitTransfer posts a transfer request between (sub) accounts
func (c *Client) SubmitTransfer(
	amt string, ccy string, dst int, result *inout.XferSubAcctOut) error {

	return c.Call(
		"private/submit_transfer_to_subaccount",
		&inout.XferSubAcctIn{Amt: amt, Ccy: ccy, Dst: dst},
		result)
}
