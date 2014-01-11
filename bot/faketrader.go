package bot

import (
	"fmt"
	"github.com/thinxer/graphpipe"

	s "github.com/thinxer/gocoins"
	a "github.com/thinxer/gotrader/analytic"
)

// TODO: delay, partial transaction, initial balance
type FakeTrader struct {
	orders  []s.Order
	balance map[s.Symbol]float64

	lastOrder int
	market    a.TickSource
	order     OrderSource
}

type FakeTraderConfig struct {
}

func newFakeTrader(config *FakeTraderConfig) (*FakeTrader, error) {
	return &FakeTrader{balance: make(map[s.Symbol]float64)}, nil
}

func (f *FakeTrader) SetInput(market a.TickSource, order OrderSource) {
	f.market = market
	f.order = order
}

func (f *FakeTrader) Update(tid int) (updated bool) {
	ordertid, order := f.order.Value()
	if ordertid > f.lastOrder {
		f.lastOrder = ordertid
		f.orders = append(f.orders, *order)
	}
	tickid, tick := f.market.Value()
	if tickid == 0 {
		return false
	}
	i := 0
	for i < len(f.orders) {
		o := &f.orders[i]
		price := o.Price

		if o.Type == s.Buy && price >= tick.Low {
			if price > tick.High {
				price = tick.High
			}
			f.balance[o.Pair.Target] += o.Amount
			f.balance[o.Pair.Base] -= o.Amount * price
			f.orders = append(f.orders[:i], f.orders[i+1:]...)
			updated = true
		} else if o.Type == s.Sell && price <= tick.High {
			if price < tick.Low {
				price = tick.Low
			}
			f.balance[o.Pair.Target] -= o.Amount
			f.balance[o.Pair.Base] += o.Amount * price
			f.orders = append(f.orders[:i], f.orders[i+1:]...)
			updated = true
		} else {
			i++
		}
	}
	if updated {
		if len(f.orders) > 0 {
			fmt.Println("orders:", f.orders)
		}
		fmt.Println("balance:", f.balance)
	}
	return
}

func (f *FakeTrader) Closed() bool {
	return f.order.Closed() || f.market.Closed()
}

func init() {
	graphpipe.Regsiter("FakeTrader", newFakeTrader)
}
