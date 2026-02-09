package app

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/examples/week_4/di/platform/pkg/closer"
	ufo_v1 "github.com/mbakhodurov/examples/week_4/di/shared/pkg/proto/ufo/v1"
	v1 "github.com/mbakhodurov/examples/week_4/di/ufo/internal/api/ufo/v1"
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/config"
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/repository"
	ufo_repo "github.com/mbakhodurov/examples/week_4/di/ufo/internal/repository/ufo"
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/service"
	ufo_service "github.com/mbakhodurov/examples/week_4/di/ufo/internal/service/ufo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type diConainer struct {
	ufoV1Api ufo_v1.UFOServiceServer

	ufoService service.UFOService
	ufoRepo    repository.UFORepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diConainer {
	return &diConainer{}
}

func (d *diConainer) UfoV1Api(ctx context.Context) ufo_v1.UFOServiceServer {
	if d.ufoV1Api == nil {
		d.ufoV1Api = v1.NewApi(d.ufoService)
	}
	return d.ufoV1Api
}

func (d *diConainer) UfoService(ctx context.Context) service.UFOService {
	if d.ufoService == nil {
		d.ufoService = ufo_service.NewService(d.ufoRepo)
	}
	return d.ufoService
}

func (d *diConainer) UfoRepo(ctx context.Context) repository.UFORepository {
	if d.ufoRepo == nil {
		d.ufoRepo = ufo_repo.NewRepository(d.mongoDBClient)
	}
	return d.ufoRepo
}

func (d *diConainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfing().Mongo.URI()))
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

		d.mongoDBClient = client
	}
	return d.mongoDBClient
}

func (d *diConainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.mongoDBClient.Database(config.AppConfing().Mongo.DatabaseName())
	}
	return d.mongoDBHandle
}
