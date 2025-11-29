package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	ufo_v1 "github.com/mbakhodurov/examples/week_1/grpc/pkg/proto/ufo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type ufoService struct {
	ufo_v1.UnimplementedUFOServiceServer

	mu        sync.RWMutex
	sightings map[string]*ufo_v1.Sighting
}

func (s *ufoService) Update(_ context.Context, req *ufo_v1.UpdateRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sighting, exists := s.sightings[req.GetUuid()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
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

func (s *ufoService) Delete(_ context.Context, req *ufo_v1.DeleteRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sighting, exists := s.sightings[req.GetUuid()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
	}

	sighting.DeletedAt = timestamppb.New(time.Now())
	log.Printf("Удалено наблюдение с UUID %s", req.GetUuid())
	return &emptypb.Empty{}, nil
}

func (s *ufoService) GetAll(_ context.Context, req *ufo_v1.GetAllRequest) (*ufo_v1.GetAllResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sightings []*ufo_v1.Sighting
	for _, sighting := range s.sightings {
		sightings = append(sightings, sighting)
	}

	return &ufo_v1.GetAllResponse{
		Sightings:  sightings,
		TotalCount: int32(len(sightings)),
	}, nil
}

func (s *ufoService) Get(_ context.Context, req *ufo_v1.GetRequest) (*ufo_v1.GetResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sighting, exists := s.sightings[req.GetUuid()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
	}
	return &ufo_v1.GetResponse{
		Sighting: sighting,
	}, nil
}

func (s *ufoService) Create(_ context.Context, req *ufo_v1.CreateRequest) (*ufo_v1.CreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newUUID := uuid.NewString()
	sightingNew := &ufo_v1.Sighting{
		Uuid:      newUUID,
		Info:      req.GetInfo(),
		CreatedAt: timestamppb.New(time.Now()),
	}

	s.sightings[newUUID] = sightingNew
	log.Printf("Создано наблюдение с UUID %s", newUUID)
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

	s := grpc.NewServer()

	service := &ufoService{
		sightings: make(map[string]*ufo_v1.Sighting),
	}

	ufo_v1.RegisterUFOServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC server listening at %v\n", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
