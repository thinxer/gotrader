package analytic

import (
	s "github.com/thinxer/gocoins"
)

type TradeSource interface {
	Value() (int, *s.Trade)
	Closed() bool
}

type Float64Source interface {
	Value() (int, float64)
	Closed() bool
}
