package analytic

import "github.com/thinxer/graphpipe"

// This thing will aggregate Trades into Ticks (OHLCs).
// Please be aware that empty ticks (volume=0) won't be outputed.
type TickFormer struct {
	tid   int
	value *Tick

	tempValue *Tick
	tempStart int64
	tempTid   int
	interval  int64
	closing   bool
	closed    bool
	source    TradeSource
}

type TickFormerConfig struct {
	Interval int
}

func newTickFormer(config *TickFormerConfig) (*TickFormer, error) {
	return &TickFormer{interval: int64(config.Interval), tempStart: -1}, nil
}

func (t *TickFormer) SetInput(source TradeSource) {
	t.source = source
}

func (v *TickFormer) Update(_ int) bool {
	tid, trade := v.source.Value()
	if v.tempStart < 0 {
		v.tempStart = trade.Timestamp
	}
	if v.source.Closed() {
		if v.closing {
			v.closed = true
			return false
		} else {
			v.closing = true
		}
	}

	// output
	updated := v.closing || v.tempStart+int64(v.interval) <= trade.Timestamp
	if updated {
		v.tid, v.value = v.tempTid, v.tempValue
		for v.tempStart+v.interval <= trade.Timestamp {
			v.tempStart += v.interval
		}
	}

	// aggregate
	if updated || v.tempValue == nil {
		lastPrice := trade.Price
		if v.tempValue != nil {
			lastPrice = v.tempValue.Close
		}
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
	v.tempTid = tid
	return updated
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
