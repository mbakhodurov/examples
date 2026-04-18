package ufo

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepo.Create(ctx, info)
	if err != nil {
		return "", err
	}

	sighting, err := s.ufoRepo.Get(ctx, uuid)
	if nil == err {
		if err := s.ufoCache.Set(ctx, uuid, sighting, s.cacheTTL); err != nil {
			fmt.Println(err)
		}
	}
	return uuid, nil
}
