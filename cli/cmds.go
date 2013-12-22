package main

import (
	"log"

	s "github.com/thinxer/gocoins"
)

func history(c s.Client) {
	trades, _, err := c.History(*flagPair, -1)
	if err != nil {
		panic(err)
	}
	for _, t := range trades {
		log.Println(t)
	}
}

func orderbook(c s.Client) {
	log.Println(c.Orderbook(*flagPair, 50))
}

func watch(c s.Client) {
	ct := make(chan s.Trade)
	go c.Stream(*flagPair, -1, ct)
	for t := range ct {
		log.Println(t)
	}
}

func init() {
	cmds["watch"] = watch
	cmds["history"] = history
	cmds["orderbook"] = orderbook
}
