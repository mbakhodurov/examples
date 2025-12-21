package repository

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
)

type UFORepository interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	GetAll(ctx context.Context) ([]model.Sighting, error)
	Delete(ctx context.Context, uuid string) error
	Update(ctx context.Context, uuid string, info model.SightingUpdateInfo) error
}
