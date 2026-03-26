package ufo

import (
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository"
	def "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepo         repository.UFORepository
	telegramService def.TelegramService
}

func NewService(ufoRepository repository.UFORepository, telegramService def.TelegramService) *service {
	return &service{
		ufoRepo:         ufoRepository,
		telegramService: telegramService,
	}
}
