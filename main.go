package main

import (
	"context"
	"log"
	"os"

	"github.com/dominikus1993/playstation-promotion-checker-bot/cmd"
	"github.com/urfave/cli/v3"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd := &cli.Command{
		Name:  "xbox-promotion-bot",
		Usage: "parse xbox game promotions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "webhooktoken",
				Usage:    "discord webhook token",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "webhookid",
				Usage:    "discord webhhook id",
				Required: true,
			},
			&cli.FloatFlag{
				Name:     "pricePromotionPercentage",
				Aliases:  []string{"ppp"},
				Usage:    "minimum promotion percentage",
				Value:    50,
				Required: false,
			},
			&cli.StringFlag{
				Name:     "mongo-connection",
				Usage:    "mongodb connection string",
				Required: true,
			},
		},
		Action: cmd.XboxGamePromotionParser,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
