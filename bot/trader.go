package bot

import (
	"log"
	"time"

	s "github.com/thinxer/coincross"
	"github.com/thinxer/graphpipe"
)

type OrderSource interface {
	Value() (int, *s.Order)
	Closed() bool
}

type Trader struct {
	tid     int
	orderid int64

	client s.Client
	source OrderSource
	pair   s.Pair
}

type TraderConfig struct {
	Exchange       string
	Apikey, Secret string
	Timeout        int
	Pair           string
}

func newTrader(config *TraderConfig) (*Trader, error) {
	timeout := time.Duration(config.Timeout) * time.Second
	transport := s.TimeoutTransport(timeout, timeout)
	client := s.New(config.Exchange, config.Apikey, config.Secret, transport)
	bal, err := client.Balance()
	if err == nil {
		log.Println(bal)
	} else {
		return nil, err
	}
	trader := &Trader{client: client}
	(&trader.pair).Set(config.Pair)
	return trader, nil
}

func (t *Trader) SetSource(source OrderSource) {
	t.source = source
}

func (t *Trader) Update(tid int) graphpipe.Result {
	_, order := t.source.Value()
	for {
		log.Printf("Trading with pair: %s, dir: %s, %f@%f", t.pair, order.Type, order.Amount, order.Price)
		orderId, err := t.client.Trade(order.Type, t.pair, order.Price, order.Amount)
		switch err := err.(type) {
		case nil:
			log.Printf("Trading with pair: %s, dir: %s, %f@%f succeeded, orderId: %d", t.pair, order.Type, order.Amount, order.Price, orderId)
			t.tid = tid
			t.orderid = orderId
			return graphpipe.Update
		case s.TradeError:
			log.Printf("Trading with pair: %s, dir: %s, %f@%f failed permanently, error: %v", t.pair, order.Type, order.Amount, order.Price, err)
			return graphpipe.Skip
		default:
			log.Printf("Trading with pair: %s, dir: %s, %f@%f failed temporarily, retrying... err: %v", t.pair, order.Type, order.Amount, order.Price, err)
			time.Sleep(time.Second)
		}
	}
}

func (t *Trader) Value() (int, int64) {
	return t.tid, t.orderid
}

func (t *Trader) Closed() bool {
	return t.source.Closed()
}

func init() {
	graphpipe.Regsiter("Trader", newTrader)
}
