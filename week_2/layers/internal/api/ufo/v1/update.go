package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/layers/pkg/proto/ufo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *UFOApi) Update(ctx context.Context, req *ufo_v1.UpdateRequest) (*emptypb.Empty, error) {
	if req.UpdateInfo == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info cannot be nil")
	}
	err := api.ufoService.UpdateSight(ctx, req.GetUuid(), converter.SightingUpdateInfoToModel(req.GetUpdateInfo()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sighting: %v", err)

	}
	return &emptypb.Empty{}, nil
}
