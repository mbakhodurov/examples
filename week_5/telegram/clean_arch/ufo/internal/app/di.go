package app

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/platform/pkg/closer"
	ufo_v1 "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/shared/pkg/proto/ufo/v1"
	v1 "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/api/ufo/v1"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/client/http"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/config"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository"
	ufoRepo "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository/ufo"
	telegramSevice "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/service/telegram"
	ufoSevice "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/service/ufo"

	telegramClient "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/client/http/telegram"

	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// Захардкоженные значения для демонстрации
	telegramBotToken = "8618137347:AAG9CJyjJmeOCXEZ4uxZXYsuLqBxEHkXXY0"
)

type diContainer struct {
	mongoDbClient *mongo.Client
	mongoDBHandle *mongo.Database

	ufoRepository repository.UFORepository

	ufoService      service.UFOService
	telegramService service.TelegramService
	telegramClient  http.TelegramClient
	telegramBot     *bot.Bot

	ufov1Api ufo_v1.UFOServiceServer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) Ufov1Api(ctx context.Context) ufo_v1.UFOServiceServer {
	if di.ufov1Api == nil {
		di.ufov1Api = v1.NewApi(di.UfoService(ctx))
	}
	return di.ufov1Api
}

func (di *diContainer) UfoService(ctx context.Context) service.UFOService {
	if di.ufoService == nil {
		di.ufoService = ufoSevice.NewService(di.UfoRepository(ctx), di.TelegramService(ctx))
	}

	return di.ufoService
}

func (di *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if di.telegramService == nil {
		di.telegramService = telegramSevice.NewService(di.TelegramClient(ctx))
	}
	return di.telegramService
}

func (di *diContainer) TelegramClient(ctx context.Context) http.TelegramClient {
	if di.telegramClient == nil {
		di.telegramClient = telegramClient.NewClient(di.TelegramBot(ctx))
	}

	return di.telegramClient
}

func (di *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if di.telegramBot == nil {
		b, err := bot.New(telegramBotToken)
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		di.telegramBot = b
	}

	return di.telegramBot
}

func (di *diContainer) UfoRepository(ctx context.Context) repository.UFORepository {
	if di.ufoRepository == nil {
		di.ufoRepository = ufoRepo.NewRepository(di.MongoDBHandle(ctx))
	}
	return di.ufoRepository
}

func (di *diContainer) MongoDbClient(ctx context.Context) *mongo.Client {
	if di.mongoDbClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		di.mongoDbClient = client
	}
	return di.mongoDbClient
}

func (di *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if di.mongoDBHandle == nil {
		di.mongoDBHandle = di.MongoDbClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return di.mongoDBHandle
}
