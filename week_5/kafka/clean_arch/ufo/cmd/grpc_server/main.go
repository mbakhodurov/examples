package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/closer"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/app"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/config"
	"go.uber.org/zap"
)

const configPath = "../deploy/compose/ufo/.env"

func main() {
	// ctx := context.Background()

	// CONFIG
	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.NewApp(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}

	// err := config.Load(configPath)
	// if err != nil {
	// 	panic(fmt.Errorf("failed to load config: %w", err))
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// // Создаем MongoDB клиент
	// // Создаем MongoDB клиент
	// mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	// if err != nil {
	// 	log.Printf("failed to connect to MongoDB: %v\n", err)
	// 	return
	// }
	// defer func() {
	// 	if cerr := mongoClient.Disconnect(context.Background()); cerr != nil {
	// 		log.Printf("failed to disconnect from MongoDB: %v\n", cerr)
	// 	}
	// }()

	// // Проверяем подключение к MongoDB
	// err = mongoClient.Ping(ctx, nil)
	// if err != nil {
	// 	log.Printf("failed to ping MongoDB: %v\n", err)
	// 	return
	// }
	// log.Println("✅ Connected to MongoDB")

	// lis, err := net.Listen("tcp", config.AppConfig().UFOGRPC.Address())
	// if err != nil {
	// 	log.Printf("failed to listen: %v\n", err)
	// 	return
	// }
	// defer func() {
	// 	if cerr := lis.Close(); cerr != nil {
	// 		log.Printf("failed to close listener: %v\n", cerr)
	// 	}
	// }()

	// db := mongoClient.Database(config.AppConfig().Mongo.DatabaseName())
	// repo := ufoRepo.NewRepository(db)

	// // Logger
	// if err := logger.Init("debug", false); err != nil {
	// 	log.Fatalf("failed to init logger: %v", err)
	// }

	// // Kafka SyncProducer
	// saramaConfig := sarama.NewConfig()
	// saramaConfig.Producer.Return.Successes = true
	// saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	// saramaConfig.Producer.Retry.Max = 5

	// syncProducer, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers(), saramaConfig)
	// if err != nil {
	// 	log.Fatalf("failed to create sarama sync producer: %v", err)
	// }
	// defer syncProducer.Close()

	// // Создаем обёртку producer
	// kafkaProducer := kafkaProcuder.NewProducer(syncProducer, "ufo-recorded", logger.Logger())

	// ufoProdService := ufoProcuder.NewService(kafkaProducer)
	// // Создаем основной сервис UFO
	// service := ufoService.NewService(repo, ufoProdService)
	// // Создаем gRPC сервер
	// // s := grpc.NewServer()
	// color := "red"
	// sound := true
	// durationSeconds := int32(4)
	// fmt.Println(service.Create(ctx, model.SightingInfo{
	// 	Location:        "test_location",
	// 	Description:     "test_description",
	// 	Color:           &color,
	// 	Sound:           &sound,
	// 	DurationSeconds: &durationSeconds,
	// }))

	// consumerGroup, err := sarama.NewConsumerGroup(
	// 	config.AppConfig().Kafka.Brokers(),
	// 	"ufo-consumer-group",
	// 	saramaConfig,
	// )
	// if err != nil {
	// 	log.Fatalf("failed to create consumer group: %v", err)
	// }
	// defer consumerGroup.Close()

}
