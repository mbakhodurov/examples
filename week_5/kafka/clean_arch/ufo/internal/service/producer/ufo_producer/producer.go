package ufo_producer

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	events_v1 "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/shared/pkg/proto/events/v1"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	ufoRecordedProducer kafka.Producer
}

func NewService(ufoRecordedProducer kafka.Producer) *service {
	return &service{
		ufoRecordedProducer: ufoRecordedProducer,
	}
}

func (p *service) ProduceUFORecorded(ctx context.Context, event model.UFORecordedEvent) error {
	var observedAt *timestamppb.Timestamp
	if event.ObservedAt != nil {
		observedAt = timestamppb.New(*event.ObservedAt)
	}

	msg := &events_v1.UFORecorded{
		ObservedAt:  observedAt,
		Location:    event.Location,
		Description: event.Description,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UFORecorded", zap.Error(err))
		return err
	}

	if err := p.ufoRecordedProducer.Send(ctx, []byte(event.UUID), payload); err != nil {
		logger.Error(ctx, "failed to publish UFORecorded", zap.Error(err))
		return err
	}
	return nil
}
