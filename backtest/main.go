package main

import (
	"encoding/csv"
	"log"
	"os"
)

func main() {
	c := csv.NewReader(os.Stdin)
	for {
		if row, err := c.Read(); err == nil {
			log.Printf("%+v", row)
		} else {
			break
		}
	}
}
