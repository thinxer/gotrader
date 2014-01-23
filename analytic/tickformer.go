package analytic

import (
	"time"

	"github.com/thinxer/graphpipe"
)

// This thing will aggregate Trades into Ticks (OHLCs).
// Please be aware that empty ticks (volume=0) won't be outputed.
type TickFormer struct {
	tid   int
	value *Tick

	tempValue *Tick
	nextTick  time.Time
	pending   []*Tick
	realtime  bool

	closing bool
	closed  bool
	source  TradeSource

	interval time.Duration
}

type TickFormerConfig struct {
	Interval int
}

func newTickFormer(config *TickFormerConfig) (*TickFormer, error) {
	return &TickFormer{interval: time.Duration(config.Interval) * time.Second}, nil
}

func (t *TickFormer) SetInput(source TradeSource, timer graphpipe.NilSource) {
	t.source = source
	// No need to set timer, as it's only used to wake up this.
}

func (v *TickFormer) Update(tid int) graphpipe.Result {
	stid, trade := v.source.Value()
	if v.closed {
		return graphpipe.Skip
	} else if v.closing {
		v.closed = true
		return graphpipe.Update
	} else if v.source.Closed() {
		v.tid, v.value = tid, v.tempValue
		v.closing = true
		return graphpipe.Update | graphpipe.More
	}

	now := time.Now()
	var tradeTime time.Time
	if trade != nil {
		tradeTime = time.Unix(trade.Timestamp, 0)
	}

	// init
	if v.nextTick.IsZero() {
		if stid == tid {
			v.nextTick = tradeTime.Add(v.interval)
			price := trade.Price
			v.tempValue = &Tick{Pair: trade.Pair, Timestamp: tradeTime, Open: price, Close: price, High: price, Low: price, Volume: trade.Amount}
		}
		return graphpipe.Skip
	}

	// output
	for (v.realtime && v.nextTick.Before(now)) ||
		(!v.realtime && stid == tid && v.nextTick.Before(tradeTime)) {

		v.pending = append(v.pending, v.tempValue)
		price := v.tempValue.Close
		v.tempValue = &Tick{Pair: trade.Pair, Timestamp: v.nextTick, Open: price, Close: price, High: price, Low: price, Volume: 0}
		v.nextTick = v.nextTick.Add(v.interval)
		if v.nextTick.After(now) {
			v.realtime = true
		}
	}
	// update tick
	if stid == tid {
		v.tempValue.Volume += trade.Amount
		if v.tempValue.High < trade.Price {
			v.tempValue.High = trade.Price
		}
		if v.tempValue.Low > trade.Price {
			v.tempValue.Low = trade.Price
		}
		v.tempValue.Close = trade.Price
	}
	// output
	if len(v.pending) > 0 {
		v.tid, v.value = tid, v.pending[0]
		v.pending = v.pending[1:]
		if len(v.pending) > 0 {
			return graphpipe.Update | graphpipe.More
		} else {
			return graphpipe.Update
		}
	} else {
		return graphpipe.Skip
	}
}

func (v *TickFormer) Value() (int, *Tick) {
	return v.tid, v.value
}

func (v *TickFormer) Closed() bool {
	return v.closed
}

func init() {
	graphpipe.Regsiter("TickFormer", newTickFormer)
}
