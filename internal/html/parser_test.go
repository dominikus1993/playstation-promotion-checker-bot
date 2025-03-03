package html

import (
	"context"
	"fmt"
	"testing"

	gotolkit "github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestParsingFirstPage(t *testing.T) {
	parser := NewPlayStationStoreHtmlParser("https://store.playstation.com/pl-pl/category/3f772501-f6f8-49b7-abac-874a88ca4897")
	result := ParsePageChannel(parser, 1)
	subject := gotolkit.ToSlice(result)
	assert.NotNil(t, subject)
	assert.NotEmpty(t, subject)
	for _, game := range subject {
		link, err := game.GetLink()
		assert.NoError(t, err)
		assert.NotEmpty(t, game.Title, fmt.Sprintf("Game of link %s", link))
		assert.NotEmpty(t, link)
		oldprice := game.FormatOldPrice()
		assert.NotEmpty(t, oldprice)
		price := game.FormatPrice()
		assert.NotEmpty(t, price)
	}
}

func ParsePageChannel(parser *PlayStationStoreHtmlParser, page int) chan data.PlaystationStoreGame {
	chann := make(chan data.PlaystationStoreGame)
	go func() {
		defer close(chann)
		parser.parsePage(context.Background(), chann, page)
	}()
	return chann
}
