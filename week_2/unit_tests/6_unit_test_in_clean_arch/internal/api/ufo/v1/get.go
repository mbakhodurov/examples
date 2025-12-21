package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

func (a *UFOApi) Get(ctx context.Context, req *ufo_v1.GetRequest) (*ufo_v1.GetResponse, error) {
	sighting, err := a.ufoService.Get(ctx, req.GetUuid())
	if err != nil {
		// if errors.Is(err, model.ErrSightingNotFound) {
		// 	return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		// }
		return nil, err
	}

	return &ufo_v1.GetResponse{
		Sighting: converter.SightingModelToProto(sighting),
	}, nil
}
