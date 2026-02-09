package v1

import (
	"context"

	ufo_v1 "github.com/mbakhodurov/examples/week_4/di/shared/pkg/proto/ufo/v1"
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/converter"
)

func (a *api) GetAll(ctx context.Context, req *ufo_v1.GetAllRequest) (*ufo_v1.GetAllResponse, error) {
	sightings, err := a.ufoService.GetAll(ctx)
	if err != nil {
		return &ufo_v1.GetAllResponse{}, err
	}

	protosightings := make([]*ufo_v1.Sighting, 0, len(sightings))
	for _, v := range sightings {
		protosightings = append(protosightings, converter.SightingToProto(v))
	}
	return &ufo_v1.GetAllResponse{
		Sighting:   protosightings,
		TotalCount: int32(len(protosightings)),
	}, nil
}
