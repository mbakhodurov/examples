package ufo

import "github.com/mbakhodurov/examples/week_2/layers/internal/repository"

type service struct {
	ufoRepository repository.UFORepository
}

func NewService(repo repository.UFORepository) *service {
	return &service{
		ufoRepository: repo,
	}
}
