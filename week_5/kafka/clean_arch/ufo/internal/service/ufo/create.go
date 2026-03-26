package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepo.Create(ctx, info)
	if err != nil {
		return "", err
	}

	err = s.ufoProducerService.ProduceUFORecorded(ctx, model.UFORecordedEvent{
		UUID:        uuid,
		ObservedAt:  info.ObservedAt,
		Location:    info.Location,
		Description: info.Description,
	})
	if err != nil {
		return "", err
	}
	return uuid, nil
}
