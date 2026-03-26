package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	err := s.ufoRepo.Update(ctx, uuid, updateInfo)
	if err != nil {
		return err
	}

	return nil
}
