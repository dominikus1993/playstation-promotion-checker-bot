package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dominikus1993/playstation-promotion-checker-bot/internal/console"
	"github.com/dominikus1993/playstation-promotion-checker-bot/internal/discord"
	"github.com/dominikus1993/playstation-promotion-checker-bot/internal/files"
	"github.com/dominikus1993/playstation-promotion-checker-bot/internal/html"
	mongodb "github.com/dominikus1993/playstation-promotion-checker-bot/internal/mongodb"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/filters"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/filters/unique"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/usecase"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/writers"
	"github.com/urfave/cli/v3"
)

const playstationStoreUrl = "https://store.playstation.com/pl-pl/category/3f772501-f6f8-49b7-abac-874a88ca4897"

func XboxGamePromotionParser(ctx context.Context, cmd *cli.Command) error {
	webhookId := cmd.String("webhookid")
	webhooktoken := cmd.String("webhooktoken")
	promotionPercentage := cmd.Float("pricePromotionPercentage")
	mongoConnection := cmd.String("mongo-connection")
	client, err := mongodb.NewClient(ctx, mongoConnection, "PlaystationGames", "promotions")
	if err != nil {
		return fmt.Errorf("%w, failed to create mongo connection", err)
	}
	defer client.Close(ctx)
	slog.InfoContext(ctx, "starting xbox game promotion parser")
	newPromotionsFilter := mongodb.NewDatabaseOldPromotionFilter(client)
	mongoPromotionsWriter := mongodb.NewMongoGameWriter(client)

	fileFilter, err := files.NewTxtFileFilter(files.NewFileGameThatIWantProvider("./games.txt"))
	if err != nil {
		return fmt.Errorf("%w, failed to create file filter", err)
	}
	priceFilter := filters.NewPriceFilter(promotionPercentage)
	discord, err := discord.NewDiscordXboxGameWriter(webhookId, webhooktoken)
	if err != nil {
		return fmt.Errorf("%w, failed to create discord writer", err)
	}
	uniqueFilter := unique.NewUniqeFilter()
	consoleW := console.NewConsolePlaystationStoreGameWriter()
	broadcaster := writers.NewBroadcastPlaystationGameWriter(discord, consoleW, mongoPromotionsWriter)
	provider := html.NewPlayStationStoreHtmlParser(playstationStoreUrl)
	return usecase.NewPlaystationGamePromotionParser(provider, broadcaster, uniqueFilter, priceFilter, fileFilter, newPromotionsFilter).Parse(ctx)
}
