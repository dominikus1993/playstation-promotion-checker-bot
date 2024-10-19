package parsers

import (
	"context"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
)

type PlaystationNetworkDataParser interface {
	Provide(ctx context.Context) <-chan data.PlaystationStoreGame
}
