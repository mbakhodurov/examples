package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

func (a *UFOApi) Create(ctx context.Context, req *ufo_v1.CreateRequest) (*ufo_v1.CreateResponse, error) {
	uuid, err := a.ufoService.CreateSighting(ctx, converter.SightingInfoProtoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &ufo_v1.CreateResponse{
		Uuid: uuid,
	}, nil
}
