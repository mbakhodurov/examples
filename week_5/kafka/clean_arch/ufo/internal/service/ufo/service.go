package ufo

import (
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/repository"
	def "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepo            repository.UFORepository
	ufoProducerService def.UFOProducerService
}

func NewService(ufoRepository repository.UFORepository, ufoProducerService def.UFOProducerService) *service {
	return &service{
		ufoRepo:            ufoRepository,
		ufoProducerService: ufoProducerService,
	}
}
