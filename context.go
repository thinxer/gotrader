package main

import (
	s "github.com/thinxer/gocoins"
	"log"
	"sync"
	"time"
)

type TradingContext struct {
	Analytics *s.Analytics
	Client    s.Client
	Pair      s.Pair

	mu        sync.Mutex
	streamers []chan s.Trade
	waitGroup sync.WaitGroup
	running   bool

	balance   map[s.Symbol]float64
	orderbook *s.Orderbook
}

func MakeContext(c s.Client, p s.Pair) *TradingContext {
	return &TradingContext{
		Client:    c,
		Pair:      p,
		Analytics: s.MakeAnalytics(),
	}
}

func (tc *TradingContext) Run() {
	log.Printf("Trading Context started")
	// the increments must match the number of jobs
	tc.waitGroup.Add(3)
	go tc.stream()
	go tc.balanceRefresher()
	go tc.orderbookRefresher()
	tc.waitGroup.Wait()
}

func (tc *TradingContext) StreamTo(trades chan s.Trade) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.streamers = append(tc.streamers, trades)
}

func (tc *TradingContext) Balance() map[s.Symbol]float64 {
	return tc.balance
}

func (tc *TradingContext) Trade(tradeType s.TradeType, price, amount float64) (int64, error) {
	for {
		log.Printf("Trading with pair: %s, dir: %s, %f@%f", tc.Pair.String(), tradeType.String(), amount, price)
		orderId, err := tc.Client.Trade(tradeType, tc.Pair, price, amount)
		switch err := err.(type) {
		case nil:
			log.Printf("Trading with pair: %s, dir: %s, %f@%f succeeded", tc.Pair.String(), tradeType.String(), amount, price)
			return orderId, err
		case s.TradeError:
			log.Printf("Trading with pair: %s, dir: %s, %f@%f failed permanently", tc.Pair.String(), tradeType.String(), amount, price)
			return -1, err
		default:
			log.Printf("Trading with pair: %s, dir: %s, %f@%f failed temporarily, retrying...", tc.Pair.String(), tradeType.String(), amount, price)
			time.Sleep(time.Second)
		}
	}
}

func (tc *TradingContext) stream() {
	defer tc.waitGroup.Done()

	trades := make(chan s.Trade)
	go s.Tail(tc.Client, tc.Pair, time.Second, false, trades)
	for t := range trades {
		tc.Analytics.Stat(t)
		tc.mu.Lock()
		for _, b := range tc.streamers {
			b <- t
		}
		tc.mu.Unlock()
	}
	log.Printf("Tailing closed")
	tc.mu.Lock()
	for _, b := range tc.streamers {
		close(b)
	}
	tc.streamers = nil
	tc.mu.Unlock()
}

func (tc *TradingContext) balanceRefresher() {
	defer tc.waitGroup.Done()
	for tc.running {
		b, err := tc.Client.Balance()
		if err == nil {
			tc.balance = b.Available
		}
		time.Sleep(time.Second * 2)
	}
}

func (tc *TradingContext) orderbookRefresher() {
	defer tc.waitGroup.Done()
	for tc.running {
		tc.orderbook, _ = tc.Client.Orderbook(tc.Pair, 100)
		time.Sleep(time.Second * 2)
	}
}
