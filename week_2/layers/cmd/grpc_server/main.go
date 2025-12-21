package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	ufoV1API "github.com/mbakhodurov/examples/week_2/layers/internal/api/ufo/v1"
	"github.com/mbakhodurov/examples/week_2/layers/internal/interceptor"
	UFORepo "github.com/mbakhodurov/examples/week_2/layers/internal/repository/ufo"
	UFOSerivce "github.com/mbakhodurov/examples/week_2/layers/internal/service/ufo"

	ufo_v1 "github.com/mbakhodurov/examples/week_2/layers/pkg/proto/ufo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

type ufoService struct {
	ufo_v1.UnimplementedUFOServiceServer
	mu sync.RWMutex

	sightings map[string]*ufo_v1.Sighting
}

func (s *ufoService) Update(_ context.Context, req *ufo_v1.UpdateRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sighting, ok := s.sightings[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with uuid %s not found", req.GetUuid())
	}

	update := req.GetUpdateInfo()
	if update == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info is required")
	}

	if req.GetUpdateInfo().ObservedAt != nil {
		sighting.Info.ObservedAt = req.GetUpdateInfo().ObservedAt
	}

	if req.GetUpdateInfo().Location != nil {
		sighting.Info.Location = req.GetUpdateInfo().Location.Value
	}

	if req.GetUpdateInfo().Description != nil {
		sighting.Info.Description = req.GetUpdateInfo().Description.Value
	}

	if req.GetUpdateInfo().Color != nil {
		sighting.Info.Color = req.GetUpdateInfo().Color
	}

	if req.GetUpdateInfo().Sound != nil {
		sighting.Info.Sound = req.GetUpdateInfo().Sound
	}

	if req.GetUpdateInfo().DurationSeconds != nil {
		sighting.Info.DurationSeconds = req.GetUpdateInfo().DurationSeconds
	}

	sighting.UpdatedAt = timestamppb.New(time.Now())

	return &emptypb.Empty{}, nil
}

func (s *ufoService) GetAll(_ context.Context, req *ufo_v1.GetAllRequest) (*ufo_v1.GetAllResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sighting := make([]*ufo_v1.Sighting, 0, len(s.sightings))
	for _, s := range s.sightings {
		sighting = append(sighting, s)
	}
	return &ufo_v1.GetAllResponse{
		Sightings:  sighting,
		TotalCount: int32(len(sighting)),
	}, nil
}

func (s *ufoService) Delete(_ context.Context, req *ufo_v1.DeleteRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sightings, ok := s.sightings[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
	}
	sightings.DeletedAt = timestamppb.New(time.Now())

	return &emptypb.Empty{}, nil
}

func (s *ufoService) Get(_ context.Context, req *ufo_v1.GetRequest) (*ufo_v1.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sightings, ok := s.sightings[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
	}
	return &ufo_v1.GetResponse{
		Sighting: sightings,
	}, nil
}

func (s *ufoService) Create(_ context.Context, req *ufo_v1.CreateRequest) (*ufo_v1.CreateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.GetInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}

	newUUID := uuid.NewString()
	log.Printf("Создано наблюдение с UUID %s", newUUID)

	sighting := &ufo_v1.Sighting{
		Uuid:      newUUID,
		Info:      req.GetInfo(),
		CreatedAt: timestamppb.New(time.Now()),
	}
	s.sightings[newUUID] = sighting
	return &ufo_v1.CreateResponse{
		Uuid: newUUID,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor()),
			grpc.UnaryServerInterceptor(interceptor.LoggerInterceptor2()),
		),
	)

	repo := UFORepo.NewRepository()
	service := UFOSerivce.NewService(repo)
	api := ufoV1API.NewUFOApi(service)

	// service := &ufoService{
	// 	sightings: make(map[string]*ufo_v1.Sighting),
	// }

	ufo_v1.RegisterUFOServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	var gwServer *http.Server

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err = ufo_v1.RegisterUFOServiceHandlerFromEndpoint(
			ctx, mux, fmt.Sprintf("localhost:%d", grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("Failed to register gateway: %v\n", err)
			return
		}
		gwServer = &http.Server{
			Addr:        fmt.Sprintf(":%d", httpPort),
			Handler:     mux,
			ReadTimeout: 10 * time.Second,
		}
		log.Printf("🌐 HTTP server with gRPC-Gateway listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
