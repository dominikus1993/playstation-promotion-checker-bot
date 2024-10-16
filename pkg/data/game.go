package data

import (
	"errors"
	"net/url"
	"strings"

	"github.com/dominikus1993/go-toolkit/crypto"
	"github.com/dustin/go-humanize"
)

const baseMicrosoftPath = "https://www.microsoft.com/"

var errUrlIsEmptyErr = errors.New("uri is empty")

type PromotionPrice = float64
type RegularPrice = float64
type Title = string
type Link = string
type GameID = string

type PlaystationStoreGame struct {
	ID       GameID
	Title    Title
	link     Link
	price    PromotionPrice
	oldPrice RegularPrice
}

func joinPath(uri string) (string, error) {
	if uri == "" {
		return "", errUrlIsEmptyErr
	}
	if strings.HasPrefix(uri, baseMicrosoftPath) {
		return uri, nil
	}
	return url.JoinPath("https://www.microsoft.com/", uri) // TODO: use playsation store base path
}

func (g *PlaystationStoreGame) GetLink() (*url.URL, error) {

	uri, err := joinPath(g.link)
	if err != nil {
		return nil, err
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (g *PlaystationStoreGame) CalculatePromotionPercentage() float64 {
	return 100 - (g.price / g.oldPrice * 100)
}

func (g *PlaystationStoreGame) FormatPromotionPercentage() string {
	percentage := 100 - (g.price / g.oldPrice * 100)
	return formatPrice(percentage)
}

func (g *PlaystationStoreGame) FormatPrice() string {
	return formatPrice(g.price)
}

func (g *PlaystationStoreGame) FormatOldPrice() string {
	return formatPrice(g.oldPrice)
}

func NewXboxStoreGame(title Title, link Link, price PromotionPrice, oldPrice RegularPrice) PlaystationStoreGame {
	id, _ := crypto.GenerateId(title)
	return PlaystationStoreGame{
		ID:       id,
		Title:    title,
		link:     link,
		price:    price,
		oldPrice: oldPrice,
	}
}

func (g *PlaystationStoreGame) IsValidGame() bool {
	_, err := g.GetLink()
	if err != nil {
		return false
	}
	if g.Title == "" {
		return false
	}
	return true
}

func formatPrice(price float64) string {
	return humanize.FormatFloat("#,###.##", price)
}
