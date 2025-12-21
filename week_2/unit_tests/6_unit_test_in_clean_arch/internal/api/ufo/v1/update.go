package v1

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/converter"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *UFOApi) Update(ctx context.Context, req *ufo_v1.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetUpdateInfo() == nil {
		return nil, status.Error(codes.InvalidArgument, "update_info cannot be nil")
	}

	err := a.ufoService.UpdateSight(ctx, req.GetUuid(), converter.SightingUpdateInfoToModel(req.GetUpdateInfo()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update sighting: %v", err)

	}
	return &emptypb.Empty{}, nil
}
