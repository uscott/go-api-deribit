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
		if math.Abs(quote.Prc-ask[0] > api.SMALL || math.Abs(quote.Amt-ask[1]) > api.SMALL) {
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
	offset := 0
	for i, bid := range bidsOld {
		for _, index := range indices {
			if i == index {
				offset++
				indices = indices[1:]
				break
			}
		}
		quote := bk.Bids[i+offset]
		if math.Abs(quote.Amt-bid.Amt) > api.SMALL || math.Abs(quote.Prc-bid.Prc) > api.SMALL {
			t.Fatal("unexpected difference")
		}
	}
	orders = make([]inout.Order, 0, depth)
	nasks := len(bk.Asks)
	asksOrig := make([]api.Quote, nasks)
	indices = []int{1, 3, 4, 8}
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
	offset = 0
	for i, ask := range asksOrig {
		for _, index := range indices {
			if i == index {
				offset++
				indices = indices[1:]
				break
			}
		}
		quote := bk.Asks[i+offset]
		if math.Abs(quote.Amt-ask.Amt) > api.SMALL || math.Abs(quote.Prc-ask.Prc) > api.SMALL {
			t.Fatal("unexpected difference")
		}
	}
}
