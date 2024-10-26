package writers

import (
	"context"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"golang.org/x/sync/errgroup"
)

type PlaystationGameWriter interface {
	Write(ctx context.Context, games []data.PlaystationStoreGame) error
}

type BroadcastPlaystationGameWriter struct {
	writers []PlaystationGameWriter
}

func NewBroadcastPlaystationGameWriter(writers ...PlaystationGameWriter) *BroadcastPlaystationGameWriter {
	return &BroadcastPlaystationGameWriter{writers: writers}
}

func (writer *BroadcastPlaystationGameWriter) Write(ctx context.Context, games []data.PlaystationStoreGame) error {
	var wg errgroup.Group
	for _, notifier := range writer.writers {
		not := notifier
		wg.Go(func() error {
			return not.Write(ctx, games)
		})
	}
	return wg.Wait()
}
