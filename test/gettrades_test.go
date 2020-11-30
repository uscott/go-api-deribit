package test

import (
	"testing"
	"time"

	"github.com/uscott/go-api-deribit/api"
	"github.com/uscott/go-api-deribit/inout"
	"github.com/uscott/go-tools/tm"
)

func TestGetTrades(t *testing.T) {
	client, err := api.New(api.DfltCnfg())
	if err != nil {
		t.Fatal(err)
	}
	end := tm.UTC()
	start := end.Add(-1 * 24 * time.Hour)
	t.Logf("Start: %v\n", tm.Format0(start))
	t.Logf("End:   %v\n", tm.Format0(end))
	startStamp := int64(start.Sub(api.TimeZero) / time.Millisecond)
	endStamp := int64(end.Sub(api.TimeZero) / time.Millisecond)
	t.Logf("Start stamp: %v\n", startStamp)
	t.Logf("End stamp:   %v\n", endStamp)
	out := inout.UserTradesOut{}
	params := inout.TradesByInstrmtAndTmIn{
		Instrument:  "BTC-25DEC20",
		StartTmStmp: startStamp,
		EndTmStmp:   endStamp,
		Count:       10,
		IncludeOld:  true,
		Sorting:     "desc",
	}
	err := client.GetUserTradesByInstrumentAndTime(&params, &out)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Num trades: %d\n", len(out.Trades))
	t.Logf("Has more:   %v\n", out.HasMore)
	if !out.HasMore || len(out.Trades) == 0 {
		return
	}
	trades := make([]inout.UserTrade, len(out.Trades))
	copy(trades, out.Trades)
	for _, x := range trades {
		stamp := client.ConvertExchStmp(x.TmStmp)
		t.Logf("Trade time: %v\n", tm.Format0(stamp))
	}
	for out.HasMore {
		params.StartTmStmp = trades[0].TmStmp + 1
		if params.StartTmStmp >= params.EndTmStmp {
			break
		}
		if err = client.GetUserTradesByInstrumentAndTime(&params, &out); err != nil {
			return
		}
		buf := make([]inout.UserTrade, len(out.Trades))
		if len(out.Trades) > 0 {
			copy(buf, out.Trades)
			trades = append(buf, trades...)
			for _, x := range buf {
				stamp := client.ConvertExchStmp(x.TmStmp)
				t.Logf("Trade time: %v\n", tm.Format0(stamp))
			}
		}
	}
	params.StartTmStmp = startStamp
	params.EndTmStmp = trades[len(trades)-1].TmStmp - 1
	if params.StartTmStmp >= params.EndTmStmp {
		return
	}
	if err = client.GetUserTradesByInstrumentAndTime(&params, &out); err != nil {
		return
	}
	for len(out.Trades) > 0 {
		buf := make([]inout.UserTrade, len(out.Trades))
		copy(buf, out.Trades)
		for _, x := range buf {
			stamp := client.ConvertExchStmp(x.TmStmp)
			t.Logf("Trade time: %v\n", tm.Format0(stamp))
		}
		trades = append(trades, buf...)
		params.EndTmStmp = buf[len(buf)-1].TmStmp - 1
		if err = client.GetUserTradesByInstrumentAndTime(&params, &out); err != nil {
			return
		}
	}
	t.Log("All trade time stamps")
	for i, x := range trades {
		stamp := client.ConvertExchStmp(x.TmStmp)
		t.Logf("Trade time: %v\n", tm.Format0(stamp))
		if i >= 1 {
			t0, t1 := trades[i-1].TmStmp, trades[i].TmStmp
			if t0 < t1 || t0 == t1 && trades[i-1].TradeID <= trades[i].TradeID {
				t.Fatalf("Trade times not increasing:\n%+v\n%+v\n", trades[i-1], trades[i])
			}
		}
	}
	t.Logf("Num trades: %d\n", len(trades))
}
