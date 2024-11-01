package console

import (
	"context"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/k0kubun/pp/v3"
)

type ConsolePlaystationStoreGameWriter struct {
}

func NewConsolePlaystationStoreGameWriter() *ConsolePlaystationStoreGameWriter {
	return &ConsolePlaystationStoreGameWriter{}
}

func (w *ConsolePlaystationStoreGameWriter) Write(ctx context.Context, games []data.PlaystationStoreGame) error {
	scheme := pp.ColorScheme{
		Integer: pp.Green | pp.Bold,
		Float:   pp.Black | pp.BackgroundWhite | pp.Bold,
		String:  pp.Yellow,
	}

	// Register it for usage
	pp.SetColorScheme(scheme)
	for _, game := range games {
		pp.Println(game)
	}
	return nil
}
