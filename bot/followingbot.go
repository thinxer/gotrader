package bot

import (
	"fmt"

	s "github.com/thinxer/coincross"
	a "github.com/thinxer/gotrader/analytic"
	"github.com/thinxer/graphpipe"
)

type FollowingBot struct {
	tid   int
	value *s.Order

	level float64

	width       float64
	position    float64
	maxPosition float64
	pair        s.Pair

	source a.Float64Source
}

type FollowingBotConfig struct {
	InitialPosition, MaxPosition, Width float64
	Pair                                string
}

func newFollowingBot(config *FollowingBotConfig) (*FollowingBot, error) {
	b := &FollowingBot{
		level:       -1,
		position:    config.InitialPosition,
		maxPosition: config.MaxPosition,
		width:       config.Width,
	}
	(&b.pair).Set(config.Pair)
	return b, nil
}

func (f *FollowingBot) SetInput(source a.Float64Source) {
	f.source = source
}

func (b *FollowingBot) Update(tid int) (ret graphpipe.Result) {
	defer func() {
		if ret != graphpipe.Skip {
			fmt.Printf("FollowingBot: level=%v, pos=%v of %v\n", b.level, b.position, b.maxPosition)
		}
	}()
	if tid == 0 {
		// palce test order
		b.tid = 0
		b.value = &s.Order{
			Pair:   b.pair,
			Type:   s.Buy,
			Price:  0.01,
			Amount: 0.01,
		}
		return graphpipe.Update
	}
	_, lastPrice := b.source.Value()
	if b.level < 0 {
		b.level = lastPrice
		return graphpipe.Skip
	} else if b.position == b.maxPosition && lastPrice > b.level {
		b.level = lastPrice
	} else if b.position == 0 && lastPrice < b.level {
		b.level = lastPrice
	} else {
		if lastPrice > b.level+b.width {
			if b.position < b.maxPosition {
				amount := b.maxPosition - b.position
				b.tid = tid
				b.value = &s.Order{
					Pair:   b.pair,
					Type:   s.Buy,
					Price:  lastPrice * 1.1,
					Amount: amount,
				}
				b.position = b.position + amount
			}
			b.level = lastPrice
			return graphpipe.Update
		} else if lastPrice < b.level-b.width {
			if b.position > 0 {
				amount := b.position
				b.tid = tid
				b.value = &s.Order{
					Pair:   b.pair,
					Type:   s.Sell,
					Price:  lastPrice * 0.9,
					Amount: amount,
				}
				b.position = b.position - amount
			}
			b.level = lastPrice
			return graphpipe.Update
		}
	}
	return graphpipe.Skip
}

func (f *FollowingBot) Value() (int, *s.Order) {
	return f.tid, f.value
}

func (f *FollowingBot) Closed() bool {
	return f.source.Closed()
}

func init() {
	graphpipe.Regsiter("FollowingBot", newFollowingBot)
}
