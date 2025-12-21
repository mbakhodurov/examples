package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
)

func (s *service) UpdateSight(ctx context.Context, uuid string, info model.SightingUpdateInfo) error {
	if err := s.ufoRepository.Update(ctx, uuid, info); err != nil {
		return err
	}
	return nil
}
