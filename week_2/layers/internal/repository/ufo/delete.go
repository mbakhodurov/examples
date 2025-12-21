package ufo

import (
	"context"
	"time"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	"github.com/samber/lo"
)

func (r *repository) Delete(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sighting, ok := r.data[uuid]
	if !ok {
		return model.ErrSightingNotFound
	}
	sighting.Deleted_at = lo.ToPtr(time.Now())
	r.data[uuid] = sighting
	return nil
}
