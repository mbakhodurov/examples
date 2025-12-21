package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
)

func (s *UfoService) UpdateSight(ctx context.Context, uuid string, info model.SightingUpdateInfo) error {
	if err := s.uforepo.UpdateByUUID(ctx, uuid, info); err != nil {
		return err
	}
	return nil
}
