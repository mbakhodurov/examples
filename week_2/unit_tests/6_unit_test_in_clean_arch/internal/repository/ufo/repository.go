package ufo

import (
	"sync"

	def "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository"
	repomodel "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/model"
)

var _ def.UFORepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repomodel.Sighting
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repomodel.Sighting),
	}
}
