package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

func (a *UFOApi) GetAll(ctx context.Context, req *ufo_v1.GetAllRequest) (*ufo_v1.GetAllResponse, error) {
	sightings, err := a.ufoService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	protoSightings := make([]*ufo_v1.Sighting, 0, len(sightings))
	for _, sighting := range sightings {
		protoSightings = append(protoSightings, converter.SightingModelToProto(sighting))
	}

	return &ufo_v1.GetAllResponse{
		Sighting:   protoSightings,
		TotalCount: int32(len(protoSightings)),
	}, nil
}
