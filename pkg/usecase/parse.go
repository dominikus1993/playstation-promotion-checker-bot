package usecase

import (
	"context"
	"log/slog"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/filters"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/parsers"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/writers"
)

type PlaystationGamePromotionParser struct {
	filters  []filters.GameFilter
	provider parsers.PlaystationNetworkDataParser
	writer   writers.PlaystationGameWriter
}

func NewPlaystationGamePromotionParser(provider parsers.PlaystationNetworkDataParser, writer writers.PlaystationGameWriter, filters ...filters.GameFilter) *PlaystationGamePromotionParser {
	return &PlaystationGamePromotionParser{
		filters:  filters,
		provider: provider,
		writer:   writer,
	}
}

func (svc *PlaystationGamePromotionParser) Parse(ctx context.Context) error {
	games := svc.provider.Provide(ctx)
	result := filters.FilterPipeline(ctx, games, svc.filters...)

	gamesresult := channels.ToSlice(result)
	if len(gamesresult) == 0 {
		slog.Info("no games in promotion")
		return nil
	}

	return svc.writer.Write(ctx, gamesresult)
}
