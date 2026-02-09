package v1

import (
	"context"

	ufo_v1 "github.com/mbakhodurov/examples/week_4/config/shared/pkg/proto/ufo/v1"
	"github.com/mbakhodurov/examples/week_4/config/ufo/internal/converter"
)

func (a *api) Create(ctx context.Context, req *ufo_v1.CreateRequest) (*ufo_v1.CreateResponse, error) {
	uuid, err := a.ufoService.Create(ctx, converter.UFOInfoToModel(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &ufo_v1.CreateResponse{
		Uuid: uuid,
	}, nil
}
