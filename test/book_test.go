package test

import (
	"math"
	"testing"

	"github.com/uscott/go-api-deribit/api"
	"github.com/uscott/go-api-deribit/inout"
)

var (
	clientbt, _ = api.New(api.DfltCnfg())
)

func TestBookCreate(t *testing.T) {
	c, depth := api.BTCSWAP, 10
	var bkraw inout.BookOut
	err := clientbt.GetBook(c, depth, &bkraw)
	if err != nil {
		t.Fatal(err.Error())
	}
	bk, err := clientbt.NewBook(c, api.BTC, "future", depth)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(bkraw.Bids) != len(bk.Bids) || len(bkraw.Asks) != len(bk.Asks) {
		t.Fatal("lengths not equal")
	}
	nbids, nasks := len(bk.Bids), len(bk.Asks)
	for i, bid := range bkraw.Bids {
		quote := bk.Bids[i]
		if math.Abs(quote.Prc-bid[0]) > api.SMALL || math.Abs(quote.Amt-bid[1]) > api.SMALL {
			t.Fatal("bids not equal")
		}
	}
	for i, ask := range bkraw.Asks {
		quote := bk.Asks[i]
		if math.Abs(quote.Prc-ask[0]) > api.SMALL || math.Abs(quote.Amt-ask[1]) > api.SMALL {
			t.Fatal("asks not equal")
		}
	}
	if nbids > 0 {
		best, quote := bk.BestBid, bk.Bids[0]
		if math.Abs(best.Prc-quote.Prc) > api.SMALL || math.Abs(best.Amt-quote.Amt) > api.SMALL {
			t.Fatal("best bid is wrong")
		}
	}
	if nasks > 0 {
		best, quote := bk.BestAsk, bk.Asks[0]
		if math.Abs(best.Prc-quote.Prc) > api.SMALL || math.Abs(best.Amt-quote.Amt) > api.SMALL {
			t.Fatal("best bid is wrong")
		}
	}
}

func TestBookPrune(t *testing.T) {
	c, depth := api.BTCSWAP, 10
	bk, err := clientbt.NewBook(c, api.BTC, "future", depth)
	if err != nil {
		t.Fatal(err.Error())
	}
	orders := make([]inout.Order, 0, depth)
	nbids := len(bk.Bids)
	bidsOrig := make([]api.Quote, nbids)
	n := copy(bidsOrig, bk.Bids)
	if n < nbids {
		t.Fatalf("fewer copied than expected: %d, %d\n", n, nbids)
	}
	indices := []int{0, 3, 7}
	for _, i := range indices {
		quote := bk.Bids[i]
		ord := inout.Order{
			Amt:        quote.Amt,
			Direction:  api.DirBuy,
			Instrument: c,
			Prc:        inout.Price(quote.Prc),
		}
		orders = append(orders, ord)
	}
	if err = api.PruneOrdersFromBook(bk, orders); err != nil {
		t.Fatal(err.Error())
	}
	i, j := 0, 0
	var bid0, bid1 api.Quote
	for {
		for _, index := range indices {
			if i == index {
				i++
				indices = indices[1:]
				continue
			} else {
				break
			}
		}
		if i > nbids-1 {
			break
		}
		bid0 = bidsOrig[i]
		bid1 = bk.Bids[j]
		if math.Abs(bid1.Amt-bid0.Amt) > api.SMALL || math.Abs(bid1.Prc-bid0.Prc) > api.SMALL {
			t.Logf("index 0: %d\n", i)
			t.Logf("index 1: %d\n", j)
			t.Logf("Bid 0:  %+v\n", bid0)
			t.Logf("Bid 1:  %+v\n", bid1)
			t.Logf("Bids Orig.: %+v\n", bidsOrig)
			t.Logf("Bids:       %+v\n", bk.Bids)
			t.Fatal("unexpected difference in bids")
		}
		i++
		j++
	}
	orders = make([]inout.Order, 0, depth)
	nasks := len(bk.Asks)
	asksOrig := make([]api.Quote, nasks)
	n = copy(asksOrig, bk.Asks)
	if n < nasks {
		t.Fatalf("fewer copied than expected: %d, %d\n", n, nbids)
	}
	indices = []int{1, 3, 4, 5, 9}
	for _, i := range indices {
		quote := bk.Asks[i]
		ord := inout.Order{
			Amt:        quote.Amt,
			Direction:  api.DirSell,
			Instrument: c,
			Prc:        inout.Price(quote.Prc),
		}
		orders = append(orders, ord)
	}
	if err = api.PruneOrdersFromBook(bk, orders); err != nil {
		t.Fatal(err.Error())
	}
	i, j = 0, 0
	var ask0, ask1 api.Quote
	for {
		for _, index := range indices {
			if i == index {
				i++
				indices = indices[1:]
				continue
			} else {
				break
			}
		}
		if i > nasks-1 {
			break
		}
		ask0 = asksOrig[i]
		ask1 = bk.Asks[j]
		if math.Abs(ask1.Amt-ask0.Amt) > api.SMALL || math.Abs(ask1.Prc-ask0.Prc) > api.SMALL {
			t.Logf("index 0: %d\n", i)
			t.Logf("index 1: %d\n", j)
			t.Logf("Ask 1:  %+v\n", ask0)
			t.Logf("Ask 2:  %+v\n", ask1)
			t.Logf("Asks Orig.: %+v\n", asksOrig)
			t.Logf("Asks:       %+v\n", bk.Asks)
			t.Fatal("unexpected difference in asks")
		}
		i++
		j++
	}
}
