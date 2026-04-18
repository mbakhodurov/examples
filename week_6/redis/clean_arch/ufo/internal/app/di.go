package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/platform/pkg/cache"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/platform/pkg/cache/redis"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/platform/pkg/closer"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/platform/pkg/logger"
	ufo_v1 "github.com/mbakhodurov/examples/week_6/redis/clean_arch/shared/pkg/proto/ufo/v1"
	ufoV1API "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/api/ufo/v1"

	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/config"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/repository"
	ufoRepository "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/repository/ufo"
	ufoCacheRepository "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/repository/ufo_cache"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/service"
	ufoService "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/service/ufo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type diContainer struct {
	ufoV1API ufo_v1.UFOServiceServer

	ufoService         service.UFOService
	ufoRepo            repository.UFORepository
	ufoCacheRepository repository.UFOCacheRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	redisClient cache.RedisClient
	redisPool   *redigo.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (di *diContainer) RedisPool() *redigo.Pool {
	if di.redisPool == nil {
		di.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}
	return di.redisPool
}

func (di *diContainer) RedisClient() cache.RedisClient {
	if di.redisClient == nil {
		di.redisClient = redis.NewClient(di.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return di.redisClient
}

func (di *diContainer) UfoCacheRepository(ctx context.Context) repository.UFOCacheRepository {
	if di.ufoCacheRepository == nil {
		di.ufoCacheRepository = ufoCacheRepository.NewRepository(di.RedisClient())
	}
	return di.ufoCacheRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
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

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (di *diContainer) UfoRepo(ctx context.Context) repository.UFORepository {
	if di.ufoRepo == nil {
		di.ufoRepo = ufoRepository.NewRepository(di.MongoDBHandle(ctx))
	}
	return di.ufoRepo
}

func (d *diContainer) UfoV1API(ctx context.Context) ufo_v1.UFOServiceServer {
	if d.ufoV1API == nil {
		d.ufoV1API = ufoV1API.NewAPI(d.UfoService(ctx))
	}

	return d.ufoV1API
}

func (di *diContainer) UfoService(ctx context.Context) service.UFOService {
	if di.ufoService == nil {
		di.ufoService = ufoService.NewService(
			di.UfoRepo(ctx),
			di.UfoCacheRepository(ctx),
			config.AppConfig().Redis.CacheTTL(),
		)
	}
	return di.ufoService
}
