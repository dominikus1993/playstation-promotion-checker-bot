package mongodb

import (
	"context"

	"github.com/dominikus1993/go-toolkit/channels"
	"github.com/dominikus1993/playstation-promotion-checker-bot/pkg/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseOldPromotionFilter struct {
	client *MongoClient
}

func NewDatabaseOldPromotionFilter(client *MongoClient) *DatabaseOldPromotionFilter {
	return &DatabaseOldPromotionFilter{client: client}
}

func (f *DatabaseOldPromotionFilter) Filter(ctx context.Context, games <-chan data.PlaystationStoreGame) <-chan data.PlaystationStoreGame {
	return channels.Filter(games, filterGameThatWasInPromtoionSinceWeek(ctx, f.client), 10)
}

func filterGameThatWasInPromtoionSinceWeek(ctx context.Context, client *MongoClient) func(game data.PlaystationStoreGame) bool {
	return func(game data.PlaystationStoreGame) bool {
		col := client.GetCollection()
		opts := options.FindOne()
		res := col.FindOne(ctx, bson.D{{Key: "_id", Value: game.ID}}, opts)
		notexists := res.Err() == mongo.ErrNoDocuments
		return notexists
	}
}
