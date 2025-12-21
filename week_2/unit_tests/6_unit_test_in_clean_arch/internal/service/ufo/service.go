package ufo

import (
	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository"
	def "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/service"
)

var _ def.UFOService = (*UfoService)(nil)

type UfoService struct {
	uforepo repository.UFORepository
}

func NewUfoService(uforepo repository.UFORepository) *UfoService {
	return &UfoService{
		uforepo: uforepo,
	}
}
