package repository

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
)

type UFORepository interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	GetAll(ctx context.Context) ([]model.Sighting, error)
	Delete(ctx context.Context, uuid string) error
	Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
}
