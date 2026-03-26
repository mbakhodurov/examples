package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/closer"
	wrappedKafka "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka/producer"

	kafkaMiddleware "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/middleware/kafka"

	kafkaConverter "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/converter/kafka"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/converter/kafka/decoder"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	ufo_v1 "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/shared/pkg/proto/ufo/v1"

	v1 "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/api/ufo/v1"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/config"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/repository"
	ufoRepository "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/repository/ufo"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service"
	ufoProduceService "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service/producer/ufo_producer"
	ufoService "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service/ufo"

	ufoConsumeService "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service/consumer/ufo_consumer"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type diContainer struct {
	ufoRepository repository.UFORepository
	mongoDbClient *mongo.Client
	mongoDbHandle *mongo.Database

	ufoService service.UFOService

	syncProducer        sarama.SyncProducer
	ufoRecordedProducer wrappedKafka.Producer
	ufoProducerService  service.UFOProducerService

	consumerGroup       sarama.ConsumerGroup
	ufoRecordedConsumer wrappedKafka.Consumer
	ufoConsumerService  service.ConsumerService
	ufoRecordedDecoder  kafkaConverter.UFORecordedDecoder

	ufoV1API ufo_v1.UFOServiceServer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UFORecordedDecoder() kafkaConverter.UFORecordedDecoder {
	if d.ufoRecordedDecoder == nil {
		d.ufoRecordedDecoder = decoder.NewUFORecordedDecoder()
	}

	return d.ufoRecordedDecoder
}

func (d *diContainer) UfoConsumerService() service.ConsumerService {
	if d.ufoConsumerService == nil {
		d.ufoConsumerService = ufoConsumeService.NewService(d.UfoRecordedConsumer(), d.UFORecordedDecoder())
	}
	return d.ufoConsumerService
}

func (d *diContainer) UfoRecordedConsumer() wrappedKafka.Consumer {
	if d.ufoRecordedConsumer == nil {
		d.ufoRecordedConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().UfoRecordedConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.ufoRecordedConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().UfoRecordedConsumer.GroupID(),
			config.AppConfig().UfoRecordedProducer.Config(),
		)

		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

func (d *diContainer) MongoDbClient(ctx context.Context) *mongo.Client {
	if d.mongoDbClient == nil {
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

		d.mongoDbClient = client
	}
	return d.mongoDbClient
}

func (d *diContainer) MongoDbHandle(ctx context.Context) *mongo.Database {
	if d.mongoDbHandle == nil {
		d.mongoDbHandle = d.MongoDbClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return d.mongoDbHandle
}

func (d *diContainer) UfoRepository(ctx context.Context) repository.UFORepository {
	if d.ufoRepository == nil {
		d.ufoRepository = ufoRepository.NewRepository(d.MongoDbHandle(ctx))
	}
	return d.ufoRepository
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().UfoRecordedProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) UfoRecordedProducer(ctx context.Context) wrappedKafka.Producer {
	if d.ufoRecordedProducer == nil {
		d.ufoRecordedProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().UfoRecordedProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.ufoRecordedProducer
}

func (d *diContainer) UfoProducerService(ctx context.Context) service.UFOProducerService {
	if d.ufoProducerService == nil {
		d.ufoProducerService = ufoProduceService.NewService(d.UfoRecordedProducer(ctx))
	}
	return d.ufoProducerService
}

func (d *diContainer) UfoService(ctx context.Context) service.UFOService {
	if d.ufoService == nil {
		d.ufoService = ufoService.NewService(d.UfoRepository(ctx), d.UfoProducerService(ctx))
	}
	return d.ufoService
}

func (d *diContainer) UfoV1API(ctx context.Context) ufo_v1.UFOServiceServer {
	if d.ufoV1API == nil {
		d.ufoV1API = v1.NewApi(d.UfoService(ctx))
	}

	return d.ufoV1API
}
