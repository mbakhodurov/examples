package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
)

func (s *UfoService) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	sight, err := s.uforepo.GetByUUID(ctx, uuid)
	if err != nil {
		return model.Sighting{}, err
	}
	return sight, nil
}
