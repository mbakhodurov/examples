package main

import (
	"context"
	"log"
	"net"

	"github.com/mbakhodurov/examples/week_7/tracing/classification/internal/handler"
	"github.com/mbakhodurov/examples/week_7/tracing/classification/internal/service"
	"github.com/mbakhodurov/examples/week_7/tracing/platform/pkg/tracing"
	classification_v1 "github.com/mbakhodurov/examples/week_7/tracing/shared/pkg/proto/classification/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = ":50053"
)

type config struct{}

func (c *config) CollectorEndpoint() string { return "localhost:4317" }
func (c *config) ServiceName() string       { return "classification-service" }
func (c *config) Environment() string       { return "development" }
func (c *config) ServiceVersion() string    { return "1.0.0" }

func main() {
	ctx := context.Background()
	cfg := &config{}
	err := tracing.InitTracer(ctx, cfg)
	if err != nil {
		log.Printf("Failed to initialize tracing: %v\n", err)
		return
	}
	defer func() {
		if cerr := tracing.ShutdownTracer(ctx); cerr != nil {
			log.Printf("Failed to shutdown tracer: %v\n", cerr)
		}
	}()

	// Создание сервиса и хендлера
	classificationService := service.NewClassificationService()
	classificationHandler := handler.NewClassificationHandler(classificationService)

	// Создание gRPC сервера с трейсинг интерцептором
	server := grpc.NewServer(
		grpc.UnaryInterceptor(tracing.UnaryServerInterceptor("classification-service")),
	)
	reflection.Register(server)
	// Регистрация сервиса
	classification_v1.RegisterClassificationServiceServer(server, classificationHandler)

	// Запуск сервера
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
		return
	}

	log.Printf("Classification service listening on %s", grpcPort)
	err = server.Serve(lis)
	if err != nil {
		log.Printf("Failed to serve: %v\n", err)
	}
}
