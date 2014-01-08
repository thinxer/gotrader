package analytic

import (
	s "github.com/thinxer/gocoins"
	"github.com/thinxer/graphpipe"
)

type Tick struct {
	Pair                   s.Pair
	Timestamp              int64
	Open, Close, High, Low float64
	Volume                 float64
}

// This thing will aggeragate Trades into Ticks (OHLCs).
// Please be aware that empty ticks (volume=0) won't be outputed.
type TickFormer struct {
	tid   int
	value *Tick

	tempValue *Tick
	tempStart int64
	interval  int64
	source    TradeSource
}

type TickFormerConfig struct {
	Interval int
}

func newTickFormer(config *TickFormerConfig, source TradeSource) (*TickFormer, error) {
	return &TickFormer{source: source, interval: int64(config.Interval), tempStart: -1}, nil
}

func (v *TickFormer) Update(_ int) bool {
	tid, trade := v.source.Value()
	if v.tempStart < 0 {
		v.tempStart = trade.Timestamp
	}
	updated := false
	if v.source.Closed() || (v.tempStart+int64(v.interval) < trade.Timestamp) {
		v.value = v.tempValue
		v.tid = tid
		for v.tempStart+v.interval < trade.Timestamp {
			v.tempStart += v.interval
		}
		updated = true
	}
	lastPrice := trade.Price
	if v.tempValue != nil {
		lastPrice = v.tempValue.Close
	}
	if updated || v.tempValue == nil {
		v.tempValue = &Tick{Pair: trade.Pair, Timestamp: v.tempStart, Open: lastPrice, Close: lastPrice, High: lastPrice, Low: lastPrice, Volume: 0}
	}
	v.tempValue.Volume += trade.Amount
	if v.tempValue.High < trade.Price {
		v.tempValue.High = trade.Price
	}
	if v.tempValue.Low > trade.Price {
		v.tempValue.Low = trade.Price
	}
	v.tempValue.Close = trade.Price
	return updated
}

func (v *TickFormer) Value() (int, *Tick) {
	return v.tid, v.value
}

func (v *TickFormer) Closed() bool {
	return v.source.Closed()
}

func init() {
	graphpipe.Regsiter("TickFormer", newTickFormer)
}
