package api

import (
	"fmt"

	"github.com/uscott/go-api-deribit/inout"
	"github.com/uscott/go-tools/slice"
)

// SubscribeMode is a bool parametrizing whether a subscription
// is public or private
type SubscribeMode bool

// Enumeration of valid SubscribeMode
const (
	SubPblc         SubscribeMode = false
	SubPrvt         SubscribeMode = true
	pblcSubscribe                 = "public/subscribe"
	prvtSubscribe                 = "private/subscribe"
	pblcUnsubscribe               = "public/unsubscribe"
	prvtUnsubscribe               = "private/unsubscribe"
)

func (c *Client) rmvChannels(channels []string) {
	for _, s := range channels {
		slice.StrSlcRm(&c.subscriptions, s)
		delete(c.subscriptionsMap, s)
	}
}

// Subscribe subscribes to given channels
func (c *Client) Subscribe(channels []string, mode SubscribeMode) ([]string, error) {
	var (
		result []string
		err    error
	)
	switch mode {
	case SubPblc:
		err = c.Call(pblcSubscribe, &inout.SubscribeIn{Channels: channels}, &result)
	case SubPrvt:
		err = c.Call(prvtSubscribe, &inout.SubscribeIn{Channels: channels}, &result)
	default:
		err = fmt.Errorf("unknown Subscribe mode: %v", mode)
	}
	return result, err
}

// Unsubscribe unsubscribes from given channels
func (c *Client) Unsubscribe(channels []string, mode SubscribeMode) ([]string, error) {
	var (
		result []string
		err    error
	)
	switch mode {
	case SubPblc:
		err = c.Call(pblcUnsubscribe, &inout.SubscribeIn{Channels: channels}, &result)
	case SubPrvt:
		err = c.Call(prvtUnsubscribe, &inout.SubscribeIn{Channels: channels}, &result)
	default:
		err = fmt.Errorf("unknown Subscribe mode: %v", mode)
	}
	c.rmvChannels(channels)
	return result, err
}

// SubPrvt subscribes to given private channels
func (c *Client) SubPrvt(channels []string) ([]string, error) {
	return c.Subscribe(channels, SubPrvt)
}

// SubPblc subscribes to given public channels
func (c *Client) SubPblc(channels []string) ([]string, error) {
	return c.Subscribe(channels, SubPblc)
}

// UnsubPrvt unsubscribes from private channels
func (c *Client) UnsubPrvt(channels []string) ([]string, error) {
	return c.Unsubscribe(channels, SubPblc)
}

// UnsubPblc unsubscribes from public channels
func (c *Client) UnsubPblc(channels []string) ([]string, error) {
	return c.Unsubscribe(channels, SubPrvt)
}

// PublicSubscribe subscribes to a public channel
func (c *Client) PublicSubscribe(
	params *inout.SubscribeIn) (result inout.SubscribeOut, err error) {

	err = c.Call("public/subscribe", params, &result)
	return
}

// PublicUnsubscribe unsubscribes from a public channel
func (c *Client) PublicUnsubscribe(
	params *inout.SubscribeIn) (result inout.SubscribeOut, err error) {

	err = c.Call("public/unsubscribe", params, &result)
	c.rmvChannels(params.Channels)
	return
}

// PrivateSubscribe subscribes to a private channel
func (c *Client) PrivateSubscribe(
	params *inout.SubscribeIn) (result inout.SubscribeOut, err error) {

	err = c.Call("private/subscribe", params, &result)
	return
}

// PrivateUnsubscribe unsubscribes to a private channel
func (c *Client) PrivateUnsubscribe(
	params *inout.SubscribeIn) (result inout.SubscribeOut, err error) {

	err = c.Call("private/unsubscribe", params, &result)
	c.rmvChannels(params.Channels)
	return
}
