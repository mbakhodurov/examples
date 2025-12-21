package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	"github.com/mbakhodurov/examples/week_2/layers/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.Sighting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoSighting, ok := r.data[uuid]
	if !ok {
		return model.Sighting{}, model.ErrSightingNotFound
	}
	return converter.SightingToModel(repoSighting), nil
}
