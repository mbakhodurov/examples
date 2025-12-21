package service

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
)

type UFOService interface {
	GetSight(ctx context.Context, uuid string) (model.Sighting, error)
	CreateSight(ctx context.Context, info model.SightingInfo) (string, error)
	GetAllSights(ctx context.Context) ([]model.Sighting, error)
	DeleteSight(ctx context.Context, uuid string) error
	UpdateSight(ctx context.Context, uuid string, info model.SightingUpdateInfo) error
}
