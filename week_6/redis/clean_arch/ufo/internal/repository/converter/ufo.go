package converter

import (
	"time"

	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/model"
	repoModel "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/repository/model"
	"github.com/samber/lo"
)

func SightingInfoToRepoModel(info model.SightingInfo) repoModel.SightingInfo {
	return repoModel.SightingInfo{
		ObservedAt:      info.ObservedAt,
		Location:        info.Location,
		Description:     info.Description,
		Color:           info.Color,
		Sound:           info.Sound,
		DurationSeconds: info.DurationSeconds,
	}
}

func SightingToModel(sighting repoModel.Sighting) model.Sighting {
	return model.Sighting{
		Uuid:      sighting.Uuid,
		Info:      SightingInfoToModel(sighting.Info),
		CreatedAt: sighting.CreatedAt,
		UpdatedAt: sighting.UpdatedAt,
		DeletedAt: sighting.DeletedAt,
	}
}

func SightingInfoToModel(info repoModel.SightingInfo) model.SightingInfo {
	return model.SightingInfo{
		ObservedAt:      info.ObservedAt,
		Location:        info.Location,
		Description:     info.Description,
		Color:           info.Color,
		Sound:           info.Sound,
		DurationSeconds: info.DurationSeconds,
	}
}

func SightingToRedisView(info model.Sighting) repoModel.SightingRedisView {
	var observedAt *int64
	if info.Info.ObservedAt != nil {
		observedAt = lo.ToPtr(info.Info.ObservedAt.UnixNano())
	}

	var updatedAt *int64
	if info.UpdatedAt != nil {
		updatedAt = lo.ToPtr(info.UpdatedAt.UnixNano())
	}

	var deletedAt *int64
	if info.DeletedAt != nil {
		deletedAt = lo.ToPtr(info.DeletedAt.UnixNano())
	}

	return repoModel.SightingRedisView{
		Uuid:            info.Uuid,
		ObservedAtNs:    observedAt,
		Location:        info.Info.Location,
		Description:     info.Info.Description,
		Color:           info.Info.Color,
		Sound:           info.Info.Sound,
		DurationSeconds: info.Info.DurationSeconds,
		CreatedAtNs:     info.CreatedAt.UnixNano(),
		UpdatedAtNs:     updatedAt,
		DeletedAtNs:     deletedAt,
	}
}

// SightingFromRedisView - конвертер из Redis view в модель домена
func SightingFromRedisView(redisView repoModel.SightingRedisView) model.Sighting {
	var observedAt *time.Time
	if redisView.ObservedAtNs != nil {
		tmp := time.Unix(0, *redisView.ObservedAtNs)
		observedAt = &tmp
	}

	var updatedAt *time.Time
	if redisView.UpdatedAtNs != nil {
		tmp := time.Unix(0, *redisView.UpdatedAtNs)
		updatedAt = &tmp
	}

	var deletedAt *time.Time
	if redisView.DeletedAtNs != nil {
		tmp := time.Unix(0, *redisView.DeletedAtNs)
		deletedAt = &tmp
	}

	return model.Sighting{
		Uuid: redisView.Uuid,
		Info: model.SightingInfo{
			ObservedAt:      observedAt,
			Location:        redisView.Location,
			Description:     redisView.Description,
			Color:           redisView.Color,
			Sound:           redisView.Sound,
			DurationSeconds: redisView.DurationSeconds,
		},
		CreatedAt: time.Unix(0, redisView.CreatedAtNs),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}
