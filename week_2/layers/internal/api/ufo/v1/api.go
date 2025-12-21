package v1

import (
	"github.com/mbakhodurov/examples/week_2/layers/internal/service"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/layers/pkg/proto/ufo/v1"
)

type UFOApi struct {
	ufo_v1.UnimplementedUFOServiceServer
	ufoService service.UFOService
}

func NewUFOApi(ufoService service.UFOService) *UFOApi {
	return &UFOApi{
		ufoService: ufoService,
	}
}
