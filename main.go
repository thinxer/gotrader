package main

import (
	s "github.com/thinxer/gocoins"
	"github.com/thinxer/gocoins/btcchina"
	"time"
)

func main() {
	var client s.Client
	client = btcchina.MakeClient("", "", s.TimeoutTransport(time.Second, time.Second))
	context := MakeContext(client, s.BTC_CNY)
	context.Analytics.Add(s.MakeEMA(0.9))
	go context.Run()
	bot := NewLevelBot(0.01, 0.05, 0.1)
	go bot.OperateOn(context)
	Watch(context)
}
