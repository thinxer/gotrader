package main

import (
	"log"
	"math"

	s "github.com/thinxer/gocoins"
)

type LevelBot struct {
	level       int
	levelWidth  int
	position    float64
	maxPosition float64
	step        float64
	trend       int
}

func NewLevelBot(step, initPosition, maxPosition float64) *LevelBot {
	return &LevelBot{step: step, position: initPosition, maxPosition: maxPosition, levelWidth: 100}
}

func (b *LevelBot) OperateOn(tc *TradingContext) {
	trades := make(chan s.Trade)
	ema := s.MakeEMA(0.9)
	tc.Analytics.Add(ema)
	tc.StreamTo(trades)
	for _ = range trades {
		lastPrice := tc.Analytics.Get(ema)[0]
		curLevel := int(lastPrice) / b.levelWidth
		water := lastPrice/float64(b.levelWidth) - float64(curLevel)
		if curLevel != b.level && b.level > 0 {
			b.trend = curLevel - b.level
		}
		b.level = curLevel
		log.Printf("LevelBot: curLevel=%v, trend=%v, pos=%v, mpos=%v", curLevel, b.trend, b.position, b.maxPosition)
		if b.trend < 0 && water < 0.9 && b.position > 0 {
			amount := math.Min(b.step, b.position)
			_, err := tc.Trade(s.Sell, lastPrice-float64(b.levelWidth/10), amount)
			if err == nil {
				b.position = b.position - amount
			}
		} else if b.trend > 0 && water > 0.05 && b.position < b.maxPosition {
			amount := math.Min(b.step, b.maxPosition-b.position)
			_, err := tc.Trade(s.Buy, lastPrice+float64(b.levelWidth/10), amount)
			if err == nil {
				b.position = b.position + amount
			}
		}
	}
}
