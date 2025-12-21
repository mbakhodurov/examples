package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/converter"
)

func (r *repository) GetByUUID(ctx context.Context, uuid string) (model.Sighting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sighting, ok := r.data[uuid]
	if !ok {
		return model.Sighting{}, model.ErrSightingNotFound
	}

	return converter.SightingToModel(sighting), nil
}
