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
		log.Println(tc.Analytics.GetAll())
	}
}
