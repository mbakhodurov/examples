package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/layers/pkg/proto/ufo/v1"
)

func (api *UFOApi) Create(ctx context.Context, req *ufo_v1.CreateRequest) (*ufo_v1.CreateResponse, error) {
	uuid, err := api.ufoService.CreateSight(ctx, converter.UFOInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}
	return &ufo_v1.CreateResponse{
		Uuid: uuid,
	}, nil
}
