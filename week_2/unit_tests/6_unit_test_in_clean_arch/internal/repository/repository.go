package repository

import (
	"context"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
)

type UFORepository interface {
	Create(ctx context.Context, sighting model.SightingInfo) (string, error)
	GetByUUID(ctx context.Context, uuid string) (model.Sighting, error)
	DeleteByUUID(ctx context.Context, uuid string) error
	GetAll(ctx context.Context) ([]model.Sighting, error)
	UpdateByUUID(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
}
