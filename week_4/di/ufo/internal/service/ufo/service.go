package ufo

import (
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/repository"
	def "github.com/mbakhodurov/examples/week_4/di/ufo/internal/service"
)

type service struct {
	ufoRepo repository.UFORepository
}

var _ def.UFOService = (*service)(nil)

func NewService(ufoRepo repository.UFORepository) *service {
	return &service{
		ufoRepo: ufoRepo,
	}
}
