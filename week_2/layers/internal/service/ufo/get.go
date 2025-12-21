package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
)

func (s *service) GetSight(ctx context.Context, uuid string) (model.Sighting, error) {
	sighting, err := s.ufoRepository.Get(ctx, uuid)
	if err != nil {
		return model.Sighting{}, err
	}
	return sighting, nil
}
