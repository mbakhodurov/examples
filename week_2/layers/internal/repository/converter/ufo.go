package converter

import (
	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	repoModel "github.com/mbakhodurov/examples/week_2/layers/internal/repository/model"
)

func SightingInfoToRepoModel(info model.SightingInfo) repoModel.SightingInfo {
	return repoModel.SightingInfo{
		Observed_at:      info.Observed_at,
		Location:         info.Location,
		Description:      info.Description,
		Color:            info.Color,
		Sound:            info.Sound,
		Duration_seconds: info.Duration_seconds,
	}
}

func SightingInfoToModel(info repoModel.SightingInfo) model.SightingInfo {
	return model.SightingInfo{
		Observed_at:      info.Observed_at,
		Location:         info.Location,
		Description:      info.Description,
		Color:            info.Color,
		Sound:            info.Sound,
		Duration_seconds: info.Duration_seconds,
	}
}

func SightingToModel(sighting repoModel.Sighting) model.Sighting {
	return model.Sighting{
		Uuid:       sighting.Uuid,
		Info:       SightingInfoToModel(sighting.Info),
		Created_at: sighting.Created_at,
		Updated_at: sighting.Updated_at,
		Deleted_at: sighting.Deleted_at,
	}
}
