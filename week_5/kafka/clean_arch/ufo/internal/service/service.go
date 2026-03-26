package service

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
)

type UFOService interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
	Delete(ctx context.Context, uuid string) error
	GetAll(ctx context.Context) ([]model.Sighting, error)
}

type UFOProducerService interface {
	ProduceUFORecorded(ctx context.Context, event model.UFORecordedEvent) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
