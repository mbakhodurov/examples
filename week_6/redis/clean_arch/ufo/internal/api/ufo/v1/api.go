package v1

import (
	ufo_v1 "github.com/mbakhodurov/examples/week_6/redis/clean_arch/shared/pkg/proto/ufo/v1"
	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/service"
)

type api struct {
	ufo_v1.UnimplementedUFOServiceServer
	ufoService service.UFOService
}

func NewAPI(ufoService service.UFOService) *api {
	return &api{
		ufoService: ufoService,
	}
}
