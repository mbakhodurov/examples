package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepo.Create(ctx, info)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
