package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
)

func (s *UfoService) GetAll(ctx context.Context) ([]model.Sighting, error) {
	sightings, err := s.uforepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return sightings, nil
}
