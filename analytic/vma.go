package analytic

import (
	s "github.com/thinxer/gocoins"
	"github.com/thinxer/graphpipe"
)

type VMA struct {
	tid   int
	value float64

	trades   []*s.Trade
	sum, vol float64
	maxVol   float64
	source   TradeSource
}

type VMAConfig struct {
	Volume float64
}

func newVMA(config *VMAConfig, source TradeSource) (*VMA, error) {
	return &VMA{maxVol: config.Volume, source: source}, nil
}

func (v *VMA) Update(tid int) bool {
	tid, trade := v.source.Value()
	for v.vol > v.maxVol {
		v.vol -= v.trades[0].Amount
		v.sum -= v.trades[0].Amount * v.trades[0].Price
		v.trades = v.trades[1:]
	}
	v.trades = append(v.trades, trade)
	v.sum += trade.Amount * trade.Price
	v.vol += trade.Amount

	v.tid, v.value = tid, v.sum/v.vol
	return true
}

func (v *VMA) Value() (int, float64) {
	return v.tid, v.value
}

func (v *VMA) Closed() bool {
	return v.source.Closed()
}

func init() {
	graphpipe.Regsiter("VMA", newVMA)
}
