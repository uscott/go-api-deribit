package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/chuckpreslar/emission"
	jsoniter "github.com/json-iterator/go"
	"github.com/uscott/go-api-deribit/inout"
)

func (c *Client) subscriptionsProcess(event *Event) (*emission.Emitter, error) {
	var (
		notification interface{}
		err          error
	)
	switch {
	case event.Channel == "announcements":
		notification = &inout.NtfctnAnncmnt{}
	case strings.HasPrefix(event.Channel, "book"):
		switch strings.Count(event.Channel, ".") {
		case 2:
			// book.BTC-PERPETUAL.raw
			// book.BTC-PERPETUAL.100ms
			if strings.HasSuffix(event.Channel, ".raw") {
				notification = &inout.NtfctnOrdrBkRaw{}
			} else {
				notification = &inout.NtfctnOrdrBk{}
			}
		case 4:
			// book.BTC-PERPETUAL.none.10.100ms
			notification = &inout.NtfctnOrdrBkGrp{}
		}
	case strings.HasPrefix(event.Channel, "deribit_price_index"):
		notification = &inout.NtfctnDrbtPrcIndx{}
	case strings.HasPrefix(event.Channel, "deribit_price_ranking"):
		notification = &inout.NtfctnDrbtPrcRnk{}
	case strings.HasPrefix(event.Channel, "estimated_expiration_price"):
		notification = &inout.NtfctnEstmtdExprtnPrc{}
	case strings.HasPrefix(event.Channel, "markprice.options"):
		notification = &inout.NtfctnMrkPrcOptns{}
	case strings.HasPrefix(event.Channel, "perpetual"):
		notification = &inout.NtfctnPrptl{}
	case strings.HasPrefix(event.Channel, "quote"):
		notification = &inout.NtfctnQut{}
	case strings.HasPrefix(event.Channel, "ticker"):
		notification = &inout.NtfctnTckr{}
	case strings.HasPrefix(event.Channel, "trades"):
		notification = &inout.NtfctnTrades{}
	case strings.HasPrefix(event.Channel, "user.changes"):
		notification = &inout.NtfctnUserChgs{}
	case strings.HasPrefix(event.Channel, "user.orders"):
		if string(event.Data)[0] == '{' {
			ordr, usrOrdr := inout.Order{}, inout.NtfctnUserOrdr{}
			err := jsoniter.Unmarshal(event.Data, &ordr)
			if err != nil {
				if c.Config.DebugMode {
					c.Logger.Println(string(event.Data))
					c.Logger.Println(err.Error())
				}
				return nil, err
			}
			usrOrdr = append(usrOrdr, ordr)
			notification = &usrOrdr
		} else {
			notification = &inout.NtfctnUserOrdr{}
		}
		return c.Emit(event.Channel, notification), nil
	case strings.HasPrefix(event.Channel, "user.portfolio"):
		notification = &inout.NtfctnPrtflio{}
	case strings.HasPrefix(event.Channel, "user.trades"):
		notification = &inout.NtfctnUserTrades{}
	default:
		return nil, fmt.Errorf("%v", string(event.Data))
	}
	err = jsoniter.Unmarshal(event.Data, &notification)
	if err != nil {
		if c.Config.DebugMode {
			c.Logger.Println(string(event.Data))
			c.Logger.Println(err.Error())
		}
		return nil, err
	}
	return c.Emit(event.Channel, notification), nil
}

// TODO: Replace the huge-ass chain of if-else with a switch statement
func (c *Client) subscriptionsProcessOld(event *Event) {
	if c.Config.DebugMode {
		c.Logger.Printf("Channel: %v %v", event.Channel, string(event.Data))
	}
	if event.Channel == "announcements" {
		var notification inout.NtfctnAnncmnt
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "book") {
		count := strings.Count(event.Channel, ".")
		if count == 2 {
			// book.BTC-PERPETUAL.raw
			// book.BTC-PERPETUAL.100ms
			if strings.HasSuffix(event.Channel, ".raw") {
				var notification inout.NtfctnOrdrBkRaw
				err := jsoniter.Unmarshal(event.Data, &notification)
				if err != nil {
					log.Printf("%v", err)
					return
				}
				c.Emit(event.Channel, &notification)
			} else {
				var notification inout.NtfctnOrdrBk
				err := jsoniter.Unmarshal(event.Data, &notification)
				if err != nil {
					log.Printf("%v", err)
					return
				}
				c.Emit(event.Channel, &notification)
			}
		} else if count == 4 {
			// book.BTC-PERPETUAL.none.10.100ms
			var notification inout.NtfctnOrdrBkGrp
			err := jsoniter.Unmarshal(event.Data, &notification)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			c.Emit(event.Channel, &notification)
		}
	} else if strings.HasPrefix(event.Channel, "deribit_price_index") {
		var notification inout.NtfctnDrbtPrcIndx
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "deribit_price_ranking") {
		var notification inout.NtfctnDrbtPrcRnk
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "estimated_expiration_price") {
		var notification inout.NtfctnEstmtdExprtnPrc
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "markprice.options") {
		var notification inout.NtfctnMrkPrcOptns
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "perpetual") {
		var notification inout.NtfctnPrptl
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "quote") {
		var notification inout.NtfctnQut
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "ticker") {
		var notification inout.NtfctnTckr
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "trades") {
		var notification inout.NtfctnTrades
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.changes") {
		var notification inout.NtfctnUserChgs
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.orders") {
		if string(event.Data)[0] == '{' {
			var notification inout.NtfctnUserOrdr
			var order inout.Order
			err := jsoniter.Unmarshal(event.Data, &order)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			notification = append(notification, order)
			c.Emit(event.Channel, &notification)
		} else {
			var notification inout.NtfctnUserOrdr
			err := jsoniter.Unmarshal(event.Data, &notification)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			c.Emit(event.Channel, &notification)
		}
	} else if strings.HasPrefix(event.Channel, "user.portfolio") {
		var notification inout.NtfctnPrtflio
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.trades") {
		var notification inout.NtfctnUserTrades
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else {
		log.Printf("%v", string(event.Data))
	}
}
