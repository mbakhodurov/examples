package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
)

func (s *service) GetAllSights(ctx context.Context) ([]model.Sighting, error) {
	sightings, err := s.ufoRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return sightings, nil
}
