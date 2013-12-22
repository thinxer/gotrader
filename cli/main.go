// A simple tool for trading on the command cli
package main

import (
	"flag"
	"net/http"
	"time"

	s "github.com/thinxer/gocoins"
)

var (
	flagExchange = flag.String("exchange", "btcchina", "the exchange you would use")
	flagCmd      = flag.String("cmd", "watch", "command")
	flagPair     = &s.Pair{s.CNY, s.BTC}
	flagTimeout  = flag.Duration("timeout", 10*time.Second, "timeout for connections")

	cmds = make(map[string]func(c s.Client))
)

func init() {
	flag.Var(flagPair, "pair", "pair to operate on")
}

func main() {
	flag.Parse()

	client := s.New(*flagExchange, "", "", s.TimeoutTransport(*flagTimeout, *flagTimeout))
	cmds[*flagCmd](client)
}
