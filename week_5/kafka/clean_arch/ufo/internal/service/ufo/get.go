package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
	"go.uber.org/zap"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	sighting, err := s.ufoRepo.Get(ctx, uuid)
	if err != nil {

		logger.Error(ctx, "failed to get ufo",
			zap.String("uuid", uuid),
			zap.Error(err),
		)
		return model.Sighting{}, err
	}
	return sighting, nil
}
