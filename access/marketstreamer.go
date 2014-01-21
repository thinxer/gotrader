package access

import (
	"time"

	s "github.com/thinxer/coincross"
	"github.com/thinxer/graphpipe"
)

type MarketStreamer struct {
	tid   int
	value *s.Trade

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
	m := &MarketStreamer{trades: client.Stream(pair, int64(since)).C}
	return m, nil
}

func (m *MarketStreamer) Update(tid int) bool {
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
