// +build ignore

package analytic

import "github.com/thinxer/graphpipe"

// This thing will aggeragate Ticks into larger Ticks (OHLCs).
type TickCombiner struct {
	tid   int
	value *Tick

	tempValue *Tick
	tempStart int64
	tempTid   int
	interval  int64
	closing   bool
	closed    bool
	source    TickSource
}

type TickCombinerConfig struct {
	Interval int
}

func newTickCombiner(config *TickCombinerConfig) (*TickCombiner, error) {
	return &TickCombiner{interval: int64(config.Interval), tempStart: -1}, nil
}

func (t *TickCombiner) SetInput(source TickSource) {
	t.source = source
}

func (v *TickCombiner) Update(_ int) graphpipe.Result {
	tid, tick := v.source.Value()
	if v.tempStart < 0 {
		v.tempStart = tick.Timestamp
	}

	if v.source.Closed() {
		if v.closing {
			v.closed = true
			return graphpipe.Skip
		} else {
			v.closing = true
		}
	}

	// output
	updated := v.closing || v.tempStart+int64(v.interval) <= tick.Timestamp
	if updated {
		v.tid, v.value = v.tempTid, v.tempValue
		for v.tempStart+v.interval <= tick.Timestamp {
			v.tempStart += v.interval
		}
	}

	// aggregate
	if updated || v.tempValue == nil {
		// Make a copy of tick
		newTick := *tick
		v.tempValue = &newTick
	}
	v.tempValue.Volume += tick.Volume
	if v.tempValue.High < tick.High {
		v.tempValue.High = tick.High
	}
	if v.tempValue.Low > tick.Low {
		v.tempValue.Low = tick.Low
	}
	v.tempValue.Close = tick.Close
	v.tempTid = tid

	if updated {
		return graphpipe.Update
	}
	return graphpipe.Skip
}

func (v *TickCombiner) Value() (int, *Tick) {
	return v.tid, v.value
}

func (v *TickCombiner) Closed() bool {
	return v.closed
}

func init() {
	graphpipe.Regsiter("TickCombiner", newTickCombiner)
}
