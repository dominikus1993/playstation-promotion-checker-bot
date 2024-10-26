package filters

import (
	"context"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
)

type GameFilter interface {
	Filter(ctx context.Context, games <-chan data.PlaystationStoreGame) <-chan data.PlaystationStoreGame
}

func FilterPipeline(ctx context.Context, stream <-chan data.PlaystationStoreGame, predicates ...GameFilter) <-chan data.PlaystationStoreGame {
	out := stream
	for _, predicate := range predicates {
		tmp := out
		out = predicate.Filter(ctx, tmp)
	}
	return out
}
