package utils

import (
	"github.com/gocolly/colly"
)

func ConnectColly() *colly.Collector {
	c := colly.NewCollector()

	return c
}
