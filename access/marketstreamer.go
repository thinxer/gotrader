package access

import (
	"time"

	s "github.com/thinxer/gocoins"
	"github.com/thinxer/graphpipe"
)

type MarketStreamer struct {
	tid   int
	value *s.Trade

	since   int64
	started bool
	closed  bool
	client  s.Client
	pair    s.Pair
	trades  chan s.Trade
}

type MarketStreamerConfig struct {
	Exchange string
	Pair     string
	Timeout  int
	Since    int64
}

func newMarketStreamer(config *MarketStreamerConfig, last int) (*MarketStreamer, error) {
	timeout := time.Duration(config.Timeout) * time.Second
	client := s.New(config.Exchange, "", "", s.TimeoutTransport(timeout, timeout))
	ms := &MarketStreamer{client: client, trades: make(chan s.Trade), since: int64(last)}
	(&ms.pair).Set(config.Pair)
	return ms, nil
}

func (m *MarketStreamer) Update(tid int) bool {
	if !m.started {
		m.started = true
		go m.client.Stream(m.pair, m.since, m.trades)
	}
	trade, ok := <-m.trades
	if ok {
		m.tid, m.value = tid, &trade
		return true
	} else {
		m.closed = true
		return false
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
