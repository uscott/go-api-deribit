package test

import (
	"testing"

	"github.com/uscott/go-api-deribit/api"
)

var (
	client, _ = api.New(api.DfltCnfg())
)

func TestBookCreate(t *testing.T) {
	c := api.BTCSWAP
	bk, err := client.NewBook(c, api.BTC, "future", 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Book: %+v\n", *bk)
}
