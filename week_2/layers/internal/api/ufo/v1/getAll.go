package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/layers/pkg/proto/ufo/v1"
)

func (api *UFOApi) GetAll(ctx context.Context, req *ufo_v1.GetAllRequest) (*ufo_v1.GetAllResponse, error) {
	sightings, err := api.ufoService.GetAllSights(ctx)
	if err != nil {
		return nil, err
	}

	protoSightings := make([]*ufo_v1.Sighting, 0, len(sightings))
	for _, sighting := range sightings {
		protoSightings = append(protoSightings, converter.SightingToProto(sighting))
	}

	return &ufo_v1.GetAllResponse{
		Sightings:  protoSightings,
		TotalCount: int32(len(protoSightings)),
	}, nil
}
