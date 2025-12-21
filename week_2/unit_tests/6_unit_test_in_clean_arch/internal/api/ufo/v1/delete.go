package v1

import (
	"context"
	"errors"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *UFOApi) Delete(ctx context.Context, req *ufo_v1.DeleteRequest) (*emptypb.Empty, error) {
	err := a.ufoService.DeleteSight(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrSightingNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
