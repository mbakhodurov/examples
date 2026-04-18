package ufo

import (
	"context"

	"github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/model"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	// Обновляем в MongoDB
	err := s.ufoRepo.Update(ctx, uuid, updateInfo)
	if err != nil {
		return err
	}

	// Инвалидируем кеш (игнорируем ошибки)
	_ = s.ufoCache.Delete(ctx, uuid)

	return nil
}
