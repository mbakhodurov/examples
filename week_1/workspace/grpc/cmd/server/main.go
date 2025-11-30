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
	ufo_v1 "github.com/mbakhodurov/examples/week_1/grpc_with_interceptor/pkg/proto/ufo/v1"
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

	mu       sync.RWMutex
	sighting map[string]*ufo_v1.Sighting
}

func NewUfoService() *ufoService {
	return &ufoService{
		sighting: make(map[string]*ufo_v1.Sighting),
	}
}

func (u *ufoService) Create(_ context.Context, rq *ufo_v1.CreateRequest) (*ufo_v1.CreateResponse, error) {
	if err := rq.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}
	u.mu.Lock()
	defer u.mu.Unlock()

	if rq.GetInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "info is required")
	}

	newUUID := uuid.NewString()

	sighting := &ufo_v1.Sighting{
		Uuid:      newUUID,
		Info:      rq.GetInfo(),
		CreatedAt: timestamppb.New(time.Now()),
	}

	// fmt.Println(sighting)
	u.sighting[newUUID] = sighting
	log.Printf("Создано наблюдение с UUID %s", newUUID)

	return &ufo_v1.CreateResponse{
		Uuid: newUUID,
	}, nil
}

func (u *ufoService) Delete(_ context.Context, req *ufo_v1.DeleteRequest) (*emptypb.Empty, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	sighting, ok := u.sighting[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with uuid %s not found", req.GetUuid())
	}
	sighting.DeletedAt = timestamppb.New(time.Now())
	return &emptypb.Empty{}, nil
}

func (u *ufoService) GetAll(_ context.Context, req *ufo_v1.GetAllRequest) (*ufo_v1.GetAllResponse, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	sightings := make([]*ufo_v1.Sighting, 0, len(u.sighting))
	for _, s := range u.sighting {
		sightings = append(sightings, s)
	}
	return &ufo_v1.GetAllResponse{
		Sightings:  sightings,
		TotalCount: int32(len(sightings)),
	}, nil
}

func (u *ufoService) Get(ctx context.Context, req *ufo_v1.GetRequest) (*ufo_v1.GetResponse, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	sighting, ok := u.sighting[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
	}
	return &ufo_v1.GetResponse{
		Sighting: sighting,
	}, nil
}

func (u *ufoService) Update(_ context.Context, req *ufo_v1.UpdateRequest) (*emptypb.Empty, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	sighting, ok := u.sighting[req.GetUuid()]
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
	s := grpc.NewServer()
	service := &ufoService{
		sighting: make(map[string]*ufo_v1.Sighting),
	}

	ufo_v1.RegisterUFOServiceServer(s, service)

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
