package mongodb

import (
	"time"

	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type screapedPlaystationGame struct {
	ID        string             `bson:"_id,omitempty"`
	Title     string             `bson:"Title"`
	Link      string             `bson:"Link"`
	Price     string             `bson:"Price"`
	OldPrice  string             `bson:"OldPrice"`
	Promotion string             `bson:"Promotion"`
	CrawledAt primitive.DateTime `bson:"CrawledAt"`
}

func fromXboxGame(game data.PlaystationStoreGame) screapedPlaystationGame {
	link, err := game.GetLink()

	if err != nil {
		panic(err)
	}

	return screapedPlaystationGame{
		ID:        game.ID,
		Title:     game.Title,
		Link:      link.String(),
		Price:     game.FormatPrice(),
		Promotion: game.FormatPromotionPercentage(),
		OldPrice:  game.FormatOldPrice(),
		CrawledAt: primitive.NewDateTimeFromTime(time.Now()),
	}
}

func toMongoWriteModel(games []data.PlaystationStoreGame) []mongo.WriteModel {
	return lo.Map(games, func(game data.PlaystationStoreGame, _ int) mongo.WriteModel {
		return mongo.NewInsertOneModel().SetDocument(fromXboxGame(game))
	})
}
