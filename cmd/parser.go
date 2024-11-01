package cmd

import (
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
	"github.com/urfave/cli/v2"
)

const playstationStoreUrl = "https://store.playstation.com/pl-pl/category/83a687fe-bed7-448c-909f-310e74a71b39"

func XboxGamePromotionParser(context *cli.Context) error {
	webhookId := context.String("webhookid")
	webhooktoken := context.String("webhooktoken")
	promotionPercentage := context.Float64("pricePromotionPercentage")
	mongoConnection := context.String("mongo-connection")
	client, err := mongodb.NewClient(context.Context, mongoConnection, "Games", "promotions")
	if err != nil {
		return fmt.Errorf("%w, failed to create mongo connection", err)
	}
	defer client.Close(context.Context)
	slog.InfoContext(context.Context, "starting xbox game promotion parser")
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
	return usecase.NewPlaystationGamePromotionParser(provider, broadcaster, uniqueFilter, priceFilter, fileFilter, newPromotionsFilter).Parse(context.Context)
}
