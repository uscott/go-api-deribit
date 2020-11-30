package test

import (
	"testing"

	"github.com/uscott/go-api-deribit/api"
	"github.com/uscott/go-api-deribit/inout"
)

func TestGetStopHistory(t *testing.T) {
	client, err := api.New(api.DfltCnfg())
	if err != nil {
		t.Fatal(err)
	}
	params := inout.StopOrderHistoryIn{Ccy: api.BTC, Count: 20}
	var hist inout.StopOrderHistoryOut
	err = client.GetStopOrderHistory(&params, &hist)
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hist.Entries {
		t.Logf("%+v\n", h)
	}
	if len(hist.Entries) > 1 {
		for i := 1; i < len(hist.Entries); i++ {
			x := hist.Entries[i-1]
			y := hist.Entries[i]
			t0 := client.ConvertExchStmp(x.TmStmp)
			t1 := client.ConvertExchStmp(y.TmStmp)
			dt := t0.Sub(t1).Minutes()
			t.Logf("Minutes: %.2f\n", dt)
		}
	}
}
