package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	"github.com/mbakhodurov/examples/week_2/layers/internal/repository/converter"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Sighting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sightings := make([]model.Sighting, 0, len(r.data))
	for _, repoSighting := range r.data {
		sightings = append(sightings, converter.SightingToModel(repoSighting))
	}
	return sightings, nil
}
