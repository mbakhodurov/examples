package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
)

func (s *UfoService) CreateSighting(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.uforepo.Create(ctx, info)
	if err != nil {
		return "", err
	}
	return uuid, nil
}
