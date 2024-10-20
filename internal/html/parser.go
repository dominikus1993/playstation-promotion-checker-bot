package html

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/gocolly/colly/v2"
	"golang.org/x/sync/errgroup"
)

var decimalRegExp = regexp.MustCompile(`[^0-9.,-]`)
var emptySpace = []byte("")
var commaRegExp = regexp.MustCompile(`,`)

type PlayStationStoreHtmlParser struct {
	psnUrl string
}

func getPageUrl(baseUrl string, page int) string {
	return baseUrl + fmt.Sprintf("/%d", page)
}

func NewPlayStationStoreHtmlParser(psnUrl string) *PlayStationStoreHtmlParser {
	return &PlayStationStoreHtmlParser{
		psnUrl: psnUrl,
	}
}

func (p *PlayStationStoreHtmlParser) Provide(ctx context.Context) <-chan data.PlaystationStoreGame {
	ch := make(chan data.PlaystationStoreGame)
	go func() {
		defer close(ch)
		p.parseAllGames(ctx, ch)
	}()
	return ch
}

func (p *PlayStationStoreHtmlParser) parseAllGames(ctx context.Context, ch chan<- data.PlaystationStoreGame) {
	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < 10; i++ {
		page := i + 1
		g.Go(func() error {
			return p.parsePage(ctx, ch, page)
		})
	}
	err := g.Wait()

	if err != nil {
		fmt.Println(err)
	}
}

func (p *PlayStationStoreHtmlParser) parsePage(ctx context.Context, ch chan<- data.PlaystationStoreGame, page int) error {
	url := getPageUrl(p.psnUrl, page)
	collector := newCollyCollector()
	collector.OnHTML("div.card", func(e *colly.HTMLElement) {

	})
	collector.OnError(func(r *colly.Response, err error) {
		slog.ErrorContext(ctx, "error while parsing page", "error", err, "url", url)
	})
	err := collector.Visit(url)
	if err != nil {
		slog.ErrorContext(ctx, "error while parsing page", "error", err, "url", url)
	}
	collector.Wait()
	return nil
}
