package main

import (
	"os"
	"github.com/thinxer/graphpipe"

	// support libraries
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/thinxer/gocoins/all"

	// Register various nodes
	_ "github.com/thinxer/gotrader/analytic"
	_ "github.com/thinxer/gotrader/bot"
	_ "github.com/thinxer/gotrader/persistence"
)

func main() {
	pipe, err := graphpipe.GraphPipeFromYAMLFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	pipe.Run()
}
