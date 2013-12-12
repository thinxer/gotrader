package main

import (
	s "github.com/thinxer/gocoins"
	"github.com/thinxer/gocoins/btcchina"
	"log"
	"time"
)

func main() {
	var client s.Client
	config := Config()
	client = btcchina.NewClient(
		config.Authorization["btcchina"].ApiKey,
		config.Authorization["btcchina"].Secret,
		s.TimeoutTransport(5*time.Second, 10*time.Second),
	)
	balance, err := client.Balance()
	if err == nil {
		log.Printf("Balance: %v", balance)
	} else {
		log.Panicf("Get balance failed. %s", err.Error())
	}
	context := NewContext(client, s.BTC_CNY)
	context.Analytics.Add(s.MakeEMA(0.9))
	go context.Run()
	bot := NewLevelBot(0.01, 0.05, 0.1)
	go bot.OperateOn(context)
	Watch(context)
}
