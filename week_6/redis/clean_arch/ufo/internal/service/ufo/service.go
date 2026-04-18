package ufo

import (
	"time"

	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/repository"
	def "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/service"
)

var _ def.UFOService = (*service)(nil)

type service struct {
	ufoRepo  repository.UFORepository
	ufoCache repository.UFOCacheRepository
	cacheTTL time.Duration
}

func NewService(ufoRepo repository.UFORepository, ufoCache repository.UFOCacheRepository, cacheTTL time.Duration) *service {
	return &service{
		ufoRepo:  ufoRepo,
		ufoCache: ufoCache,
		cacheTTL: cacheTTL,
	}
}
