package v1

import (
	"context"
	"errors"

	ufo_v1 "github.com/mbakhodurov/examples/week_4/di/shared/pkg/proto/ufo/v1"
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *api) Delete(ctx context.Context, req *ufo_v1.DeleteRequest) (*emptypb.Empty, error) {
	if err := a.ufoService.Delete(ctx, req.GetUuid()); err != nil {
		if errors.Is(err, model.ErrSightingNotFound) {
			return nil, status.Errorf(codes.NotFound, "sighting with UUID %s not found", req.GetUuid())
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
