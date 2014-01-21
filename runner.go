package main

import (
	"os"
	"github.com/thinxer/graphpipe"

	// support libraries
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/thinxer/coincross/all"

	// Register various nodes
	_ "github.com/thinxer/gotrader/access"
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
