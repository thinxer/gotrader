package bot

import (
	s "github.com/thinxer/gocoins"
	a "github.com/thinxer/gotrader/analytic"
	"github.com/thinxer/graphpipe"
	"log"
)

type FollowingBot struct {
	tid   int
	value *s.Order

	level float64

	width       float64
	position    float64
	maxPosition float64

	source a.Float64Source
}

type FollowingBotConfig struct {
	InitialPosition, MaxPosition, Width float64
}

func NewFollowingBot(config *FollowingBotConfig, source a.Float64Source) (*FollowingBot, error) {
	return &FollowingBot{
		level:       -1,
		position:    config.InitialPosition,
		maxPosition: config.MaxPosition,
		width:       config.Width,
		source:      source,
	}, nil
}

func (b *FollowingBot) Update(tid int) bool {
	if tid == 0 {
		// palce test order
		b.tid = 0
		b.value = &s.Order{
			Type:   s.Buy,
			Price:  0.01,
			Amount: 0.01,
		}
		return true
	}
	_, lastPrice := b.source.Value()
	if b.level < 0 {
		b.level = lastPrice
		return false
	} else if b.position == b.maxPosition && lastPrice > b.level {
		b.level = lastPrice
	} else if b.position == 0 && lastPrice < b.level {
		b.level = lastPrice
	} else {
		if lastPrice > b.level+b.width {
			if b.position < b.maxPosition {
				amount := b.maxPosition - b.position
				b.value = &s.Order{
					Type:   s.Buy,
					Price:  lastPrice * 1.1,
					Amount: amount,
				}
				b.position = b.position + amount
			}
			b.level = lastPrice
			return true
		} else if lastPrice < b.level-b.width {
			if b.position > 0 {
				amount := b.position
				b.value = &s.Order{
					Type:   s.Sell,
					Price:  lastPrice * 0.9,
					Amount: amount,
				}
				b.position = b.position - amount
			}
			b.level = lastPrice
			return true
		}
	}
	log.Printf("FollowingBot: level=%v, pos=%v of %v", b.level, b.position, b.maxPosition)
	return false
}

func (f *FollowingBot) Value() (int, *s.Order) {
	return f.tid, f.value
}

func (f *FollowingBot) Closed() bool {
	return f.source.Closed()
}

func init() {
	graphpipe.Regsiter("FollowingBot", NewFollowingBot)
}
