package test

import (
	"testing"

	"github.com/uscott/go-api-deribit/api"
)

var (
	client, _ = api.New(api.DfltCnfg())
)

func TestFuturesCreate(t *testing.T) {
	f, err := client.NewFuturesData(api.BTCSWAP, api.BTC)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Futures Data: %+v", *f)
}
