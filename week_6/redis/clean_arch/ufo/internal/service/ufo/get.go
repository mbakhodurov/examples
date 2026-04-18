package ufo

import (
	"context"
	"fmt"

	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	// Сначала пытаемся получить из кеша
	sighting, err := s.ufoCache.Get(ctx, uuid)
	if err == nil {
		return sighting, nil
	}

	// Если нет в кеше или ошибка, идем в MongoDB
	sighting, err = s.ufoRepo.Get(ctx, uuid)
	if err != nil {
		fmt.Println(err)
		return model.Sighting{}, err
	}

	// Сохраняем в кеш (игнорируем ошибки кеширования)
	_ = s.ufoCache.Set(ctx, uuid, sighting, s.cacheTTL)

	return sighting, nil
}
