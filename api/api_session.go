package api

import "github.com/uscott/go-api-deribit/inout"

func (c *Client) SetHeartbeat(params *inout.Heartbeat) (result string, err error) {
	err = c.Call("public/set_heartbeat", params, &result)
	return
}

func (c *Client) DisableHeartbeat() (result string, err error) {
	err = c.Call("public/disable_heartbeat", nil, &result)
	return
}

func (c *Client) EnableCancelOnDisconnect() (result string, err error) {
	err = c.Call("private/enable_cancel_on_disconnect", nil, &result)
	return
}

func (c *Client) DisableCancelOnDisconnect() (result string, err error) {
	err = c.Call("private/disable_cancel_on_disconnect", nil, &result)
	return
}
