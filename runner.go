package main

import (
	"github.com/thinxer/graphpipe"
	"os"

	_ "github.com/thinxer/gocoins/btcchina"
	_ "github.com/thinxer/gocoins/btce"

	_ "github.com/thinxer/gotrader/analytic"
	_ "github.com/thinxer/gotrader/bot"
)

func main() {
	pipe, err := graphpipe.GraphPipeFromYAMLFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	pipe.Run()
}
