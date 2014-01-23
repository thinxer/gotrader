package access

import (
	"time"

	s "github.com/thinxer/coincross"
	"github.com/thinxer/graphpipe"
)

type MarketStreamer struct {
	tid     int
	value   *s.Trade
	pending chan *s.Trade

	closed bool
	trades <-chan s.Trade
}

type MarketStreamerConfig struct {
	Exchange string
	Pair     string
	Timeout  int
	Since    int64
}

func newMarketStreamer(config *MarketStreamerConfig, since int) (*MarketStreamer, error) {
	timeout := time.Duration(config.Timeout) * time.Second
	client := s.New(config.Exchange, "", "", s.TimeoutTransport(timeout, timeout))
	var pair s.Pair
	(&pair).Set(config.Pair)
	m := &MarketStreamer{trades: client.Stream(pair, int64(since)).C, pending: make(chan *s.Trade, 1024)}
	return m, nil
}

func (m *MarketStreamer) Start(ch chan bool) {
	for trade := range m.trades {
		t := trade
		m.pending <- &t
		ch <- true
	}
	ch <- true
	close(ch)
	close(m.pending)
}

func (m *MarketStreamer) Update(tid int) graphpipe.Result {
	t, ok := <-m.pending
	if ok {
		m.tid, m.value = tid, t
		return graphpipe.Update
	} else {
		m.closed = true
		return graphpipe.Update
	}
}

func (m *MarketStreamer) Value() (int, *s.Trade) {
	return m.tid, m.value
}

func (m *MarketStreamer) Closed() bool {
	return m.closed
}

func init() {
	graphpipe.Regsiter("MarketStreamer", newMarketStreamer)
}
