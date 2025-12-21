package ufo

import (
	"context"
	"time"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	"github.com/samber/lo"
)

func (r *repository) Update(ctx context.Context, uuid string, info model.SightingUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sighting, exists := r.data[uuid]
	if !exists {
		return model.ErrSightingNotFound
	}
	if info.Observed_at != nil {
		sighting.Info.Observed_at = info.Observed_at
	}
	if info.Location != nil {
		sighting.Info.Location = *info.Location
	}

	if info.Description != nil {
		sighting.Info.Description = *info.Description
	}

	if info.Color != nil {
		sighting.Info.Color = info.Color
	}

	if info.Sound != nil {
		sighting.Info.Sound = info.Sound
	}

	if info.Duration_seconds != nil {
		sighting.Info.Duration_seconds = info.Duration_seconds
	}
	sighting.Updated_at = lo.ToPtr(time.Now())
	r.data[uuid] = sighting

	return nil
}
