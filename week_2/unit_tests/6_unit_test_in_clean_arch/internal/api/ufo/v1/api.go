package v1

import (
	ufoservice "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/service"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

type UFOApi struct {
	ufoService ufoservice.UFOService
	ufo_v1.UnimplementedUFOServiceServer
}

func NewUFOApi(ufoService ufoservice.UFOService) *UFOApi {
	return &UFOApi{
		ufoService: ufoService,
	}
}
