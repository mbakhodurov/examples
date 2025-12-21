package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/examples/week_2/layers/internal/converter"
	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/layers/pkg/proto/ufo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *UFOApi) Get(ctx context.Context, req *ufo_v1.GetRequest) (*ufo_v1.GetResponse, error) {
	sighting, err := api.ufoService.GetSight(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrSightingNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &ufo_v1.GetResponse{
		Sighting: converter.SightingToProto(sighting),
	}, nil
}
