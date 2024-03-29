package analytic

import (
	"time"
	s "github.com/thinxer/coincross"
)

type Tick struct {
	Pair                   s.Pair
	Timestamp              time.Time
	Open, Close, High, Low float64
	Volume                 float64
}

type TradeSource interface {
	Value() (int, *s.Trade)
	Closed() bool
}

type TickSource interface {
	Value() (int, *Tick)
	Closed() bool
}

type Float64Source interface {
	Value() (int, float64)
	Closed() bool
}
