package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ucpr/mongo-streamer/internal/mongo"
	"github.com/ucpr/mongo-streamer/pkg/log"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
)

type Tweet struct {
	ID        string `bson:"_id"`
	Text      string `bson:"text"`
	UserID    string `bson:"userId"`
	CreatedAt int64  `bson:"createdAt"`
	UpdateAt  int64  `bson:"updateAt"`
}

func NewRandomTweet() *Tweet {
	return &Tweet{
		ID:        fmt.Sprintf("tweet-%d", time.Now().UnixNano()),
		Text:      "Hello, World!",
		UserID:    fmt.Sprintf("user-%d", time.Now().UnixNano()),
		CreatedAt: time.Now().Unix(),
		UpdateAt:  time.Now().Unix(),
	}
}

type App struct {
	mcli *mongo.Client
}

func NewApp(mcli *mongo.Client) *App {
	return &App{
		mcli: mcli,
	}
}

func (a *App) Run(ctx context.Context) error {
	col := a.mcli.Collection("tweets")
	for {
		select {
		case <-ctx.Done():
			log.Info("Received signal to stop the application...")
			break
		default:
			operation := rand.Intn(4) // 0: insert, 1: insert & update, 2: upsert, 4: insert and delete
			tweet := NewRandomTweet()
			switch operation {
			case 0:
				if _, err := col.InsertOne(ctx, tweet); err != nil {
					log.Error("Failed to insert tweet", log.Ferror(err))
				}
			case 1:
				if _, err := col.InsertOne(ctx, tweet); err != nil {
					log.Error("Failed to insert tweet", log.Ferror(err))
					continue
				}
				tweet.Text = "Hello, World! Updated"
				if _, err := col.UpdateOne(ctx, tweet.ID, tweet); err != nil {
					log.Error("Failed to update tweet", log.Ferror(err))
				}
			case 2:
				opts := &options.UpdateOptions{
					Upsert: &[]bool{true}[0],
				}
				if _, err := col.UpdateOne(ctx, tweet.ID, tweet, opts); err != nil {
					log.Error("Failed to upsert tweet", log.Ferror(err))
				}
			case 3:
				if _, err := col.InsertOne(ctx, tweet); err != nil {
					log.Error("Failed to insert tweet", log.Ferror(err))
					continue
				}
				if _, err := col.DeleteOne(ctx, tweet.ID); err != nil {
					log.Error("Failed to delete tweet", log.Ferror(err))
				}
			}
		}

		select {
		case <-ctx.Done():
			log.Info("Received signal to stop the application...")
			break
		case <-time.After(time.Duration(rand.Intn(5)+1) * time.Second):
		}
	}
}

func (a *App) Close(ctx context.Context) error {
	return a.mcli.Disconnect(ctx)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	app, err := inject(ctx)
	if err != nil {
		log.Panic("Failed to inject", log.Ferror(err))
	}

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Error("Failed to run application", log.Ferror(err))
		}
	}()

	<-ctx.Done()
	tctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := app.Close(tctx); err != nil {
		log.Error("Failed to close application", log.Ferror(err))
		return
	}

	log.Info("Successfully graceful shutdown")
}
