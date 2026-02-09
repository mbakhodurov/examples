package service

import (
	"context"

	"github.com/mbakhodurov/examples/week_4/config/ufo/internal/model"
)

type UFOService interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
	Delete(ctx context.Context, uuid string) error
	GetAll(ctx context.Context) ([]model.Sighting, error)
}
