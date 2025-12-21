package service

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
)

type UFOService interface {
	CreateSighting(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	DeleteSight(ctx context.Context, uuids string) error
	GetAll(ctx context.Context) ([]model.Sighting, error)
	UpdateSight(ctx context.Context, uuid string, info model.SightingUpdateInfo) error
}
