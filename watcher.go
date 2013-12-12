package main

import (
	s "github.com/thinxer/gocoins"
	"log"
)

func Watch(tc *TradingContext) {
	trades := make(chan s.Trade)
	tc.StreamTo(trades)
	for trade := range trades {
		log.Println(trade)
		if tc.Orderbook() != nil {
			log.Printf("ask %v, buy %v", tc.Orderbook().Asks[0], tc.Orderbook().Bids[0])
		}
		log.Println(tc.Analytics.GetAll())
	}
}
