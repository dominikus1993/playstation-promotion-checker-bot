package html

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/gocolly/colly/v2"
	"golang.org/x/sync/errgroup"
)

type price struct {
	regularPrice   float64
	promotionPrice float64
}

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
	collector.OnHTML(".psw-grid-list", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, c *colly.HTMLElement) {
			title := c.DOM.Find(".psw-t-body").Text()
			link := c.DOM.Find("a").AttrOr("href", "")
			price, err := getPriceElement(c)
			if err != nil {
				slog.ErrorContext(ctx, "error while parsing price", "error", err, "url", url)
				return
			}
			ch <- data.NewPlaystationStoreGame(title, link, price.promotionPrice, price.regularPrice)
		})
		fmt.Println("found game")
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

func getPriceElement(e *colly.HTMLElement) (*price, error) {
	regularPrice, err := parsePrice(e.DOM.Find(".psw-c-t-2").Text())
	if err != nil {
		return nil, err
	}
	promotionPrice, err := parsePrice(e.DOM.Find(".psw-m-r-3").Text())
	if err != nil {
		return nil, err
	}
	return &price{
		regularPrice:   regularPrice,
		promotionPrice: promotionPrice,
	}, nil
}

func parsePrice(txt string) (float64, error) {
	if len(txt) == 0 {
		return 0.0, fmt.Errorf("price is empty")
	}
	// Zamiana przecinka na kropkę
	txt = string(commaRegExp.ReplaceAll([]byte(txt), []byte(".")))
	// Usunięcie wszystkich znaków, które nie są cyframi, kropką ani minusem
	priceStr := string(decimalRegExp.ReplaceAll([]byte(txt), emptySpace))
	priceFloat, err := strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
	if err != nil {
		return 0.0, err
	}

	return priceFloat, nil
}
