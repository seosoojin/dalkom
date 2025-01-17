package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nextlevellabs/go-wise/wise"
	"github.com/seosoojin/dalkom/internal/domain/binders"
	"github.com/seosoojin/dalkom/internal/domain/cards"
	"github.com/seosoojin/dalkom/internal/gateways/web"
	"github.com/seosoojin/dalkom/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start webserver",
	RunE: func(cmd *cobra.Command, args []string) error {
		if profile != "production" {
			err := godotenv.Load(".env")
			if err != nil {
				return err
			}
		}

		mongoURI := os.Getenv("MONGO_URI")

		client, err := mongo.Connect(cmd.Context(), options.Client().ApplyURI(mongoURI))
		if err != nil {
			return err
		}

		db := client.Database("dalkom")

		bindersRepo, err := wise.NewMongoSimpleRepository[models.Binder](db.Collection("binders"))
		if err != nil {
			return err
		}

		cardsRepo, err := wise.NewMongoSimpleRepository[models.Card](db.Collection("cards"))
		if err != nil {
			return err
		}

		server := web.NewServer("3000",
			binders.NewHandler(binders.NewService(bindersRepo)),
			cards.NewHandler(cards.NewService(cardsRepo)),
		)

		go server.Run()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		<-stop

		log.Println("Gracefully shutting down...")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
