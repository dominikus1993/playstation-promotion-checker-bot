package unique

import (
	"context"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
)

type uniqueFilter struct {
}

func NewUniqeFilter() *uniqueFilter {
	return &uniqueFilter{}
}

func (f *uniqueFilter) Filter(ctx context.Context, games <-chan data.PlaystationStoreGame) <-chan data.PlaystationStoreGame {
	return channels.UniqBy(games, func(game data.PlaystationStoreGame) data.GameID { return game.ID }, 10)
}
