package main

import (
	"encoding/csv"
	s "github.com/thinxer/gocoins"
	"github.com/thinxer/graphpipe"
	"os"
)

type CSVSourceConfig struct {
	filename string
}

type CSVSource struct {
	tid    int
	val    *s.Trade
	reader *csv.Reader
}

func (c *CSVSource) Update(tid int) bool {
}

func (c *CSVSource) Value() (int, *s.Trade) {
	return c.tid, c.val
}

func NewCSVSource(c *CSVSourceConfig) *CSVSource {
	reader := csv.NewReader(os.Open(c.filename))
	return &CSVSource{reader}
}
