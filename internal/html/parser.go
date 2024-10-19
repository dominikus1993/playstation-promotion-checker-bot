package html

import (
	"context"
	"regexp"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/gocolly/colly/v2"
)

var decimalRegExp = regexp.MustCompile(`[^0-9.,-]`)
var emptySpace = []byte("")
var commaRegExp = regexp.MustCompile(`,`)

type PlayStationStoreHtmlParser struct {
	psnUrl    string
	collector *colly.Collector
}

func NewPlayStationStoreHtmlParser(psnUrl string) *PlayStationStoreHtmlParser {
	return &PlayStationStoreHtmlParser{
		psnUrl:    psnUrl,
		collector: newCollyCollector(),
	}
}

func (p *PlayStationStoreHtmlParser) Provide(ctx context.Context) <-chan data.PlaystationStoreGame {
	ch := make(chan data.PlaystationStoreGame)
	go func() {
		defer close(ch)
		p.parseAllGames(ctx)
	}()
	return ch
}

func (p *PlayStationStoreHtmlParser) parseAllGames(ctx context.Context) {

}
