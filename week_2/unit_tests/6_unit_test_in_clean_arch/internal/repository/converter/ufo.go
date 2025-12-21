package converter

import (
	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	repomodel "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/model"
)

func SightingInfoToRepo(info model.SightingInfo) repomodel.SightingInfo {
	return repomodel.SightingInfo{
		Observed_at:      info.Observed_at,
		Location:         info.Location,
		Description:      info.Description,
		Color:            info.Color,
		Sound:            info.Sound,
		Duration_seconds: info.Duration_seconds,
	}
}

func SightingInfoToModel(info repomodel.SightingInfo) model.SightingInfo {
	return model.SightingInfo{
		Observed_at:      info.Observed_at,
		Location:         info.Location,
		Description:      info.Description,
		Color:            info.Color,
		Sound:            info.Sound,
		Duration_seconds: info.Duration_seconds,
	}
}

func SightingToModel(sighting repomodel.Sighting) model.Sighting {
	return model.Sighting{
		UUID:       sighting.UUID,
		Info:       SightingInfoToModel(sighting.Info),
		Created_at: sighting.Created_at,
		Updated_at: sighting.Updated_at,
		Deleted_at: sighting.Deleted_at,
	}
}

func SightingUpdateInfoToRepoModel(updateInfo model.SightingUpdateInfo) repomodel.SightingUpdateInfo {
	return repomodel.SightingUpdateInfo{
		Observed_at:      updateInfo.Observed_at,
		Location:         updateInfo.Color,
		Description:      updateInfo.Description,
		Color:            updateInfo.Color,
		Sound:            updateInfo.Sound,
		Duration_seconds: updateInfo.Duration_seconds,
	}
}
