package filters

import (
	"context"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
)

type PriceFilter struct {
	promotionPercentage float64
}

func NewPriceFilter(promotionPercentage float64) *PriceFilter {
	return &PriceFilter{promotionPercentage: promotionPercentage}
}

func (f *PriceFilter) Filter(ctx context.Context, games <-chan data.PlaystationStoreGame) <-chan data.PlaystationStoreGame {
	return channels.Filter(games, func(game data.PlaystationStoreGame) bool {
		return pricePredicate(game, f.promotionPercentage)
	}, 10)
}

func pricePredicate(game data.PlaystationStoreGame, promotionPercentage float64) bool {
	return game.CalculatePromotionPercentage() >= promotionPercentage
}
