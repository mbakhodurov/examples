package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/platform/pkg/closer"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/platform/pkg/logger"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/app"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/config"
	"go.uber.org/zap"
)

const configPath = "../deploy/compose/ufo/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

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
	// logger.Init("DEBUG", false)

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

	// db := mongoClient.Database("ufo")
	// repo := repo.NewRepository(db)

	// tgBot, err := bot.New("8618137347:AAG9CJyjJmeOCXEZ4uxZXYsuLqBxEHkXXY0")
	// if err != nil {
	// 	log.Fatalf("failed to create telegram bot: %v", err)
	// }

	// // Создаем gRPC сервер
	// s := grpc.NewServer()

	// httpTelegramClient := httpTelegramService.NewClient(tgBot)
	// telegramService := telegramService.NewService(httpTelegramClient)
	// service := ufoService.NewService(repo, telegramService)
	// api := v1.NewApi(service)

	// ufo_v1.RegisterUFOServiceServer(s, api)
	// reflection.Register(s)

	// go func() {
	// 	log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().UFOGRPC.Address())
	// 	err = s.Serve(lis)
	// 	if err != nil {
	// 		log.Printf("failed to serve: %v\n", err)
	// 		return
	// 	}
	// }()

	// // Graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("🛑 Shutting down gRPC server...")
	// s.GracefulStop()
	// log.Println("✅ Server stopped")
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
