package ufo

import (
	"context"
	"time"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	"github.com/samber/lo"
)

func (r *repository) UpdateByUUID(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sighting, ok := r.data[uuid]
	if !ok {
		return model.ErrSightingNotFound
	}
	if updateInfo.Observed_at != nil {
		sighting.Info.Observed_at = updateInfo.Observed_at
	}
	if updateInfo.Location != nil {
		sighting.Info.Location = *updateInfo.Location
	}

	if updateInfo.Description != nil {
		sighting.Info.Description = *updateInfo.Description
	}

	if updateInfo.Color != nil {
		sighting.Info.Color = updateInfo.Color
	}

	if updateInfo.Sound != nil {
		sighting.Info.Sound = updateInfo.Sound
	}

	if updateInfo.Duration_seconds != nil {
		sighting.Info.Duration_seconds = updateInfo.Duration_seconds
	}
	sighting.Updated_at = lo.ToPtr(time.Now())
	r.data[uuid] = sighting

	return nil
}
