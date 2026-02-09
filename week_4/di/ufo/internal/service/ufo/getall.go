package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/model"
)

func (s *service) GetAll(ctx context.Context) ([]model.Sighting, error) {
	sightings, err := s.ufoRepo.GetAll(ctx)
	if err != nil {
		return []model.Sighting{}, err
	}
	return sightings, nil
}
