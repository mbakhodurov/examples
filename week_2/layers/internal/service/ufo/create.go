package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
)

func (s *service) CreateSight(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
