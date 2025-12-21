package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	repoConverter "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/converter"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Sighting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sightings := make([]model.Sighting, 0, len(r.data))
	for _, sighting := range r.data {
		sightings = append(sightings, repoConverter.SightingToModel(sighting))
	}

	return sightings, nil
}
