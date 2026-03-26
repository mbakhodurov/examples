package ufo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
	repoconverter "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/repository/converter"
	repomodel "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	newUUID := uuid.NewString()

	sighting := repomodel.Sighting{
		Uuid:      newUUID,
		Info:      repoconverter.SightingInfoToRepoModel(info),
		CreatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, sighting)
	if err != nil {
		return "", err
	}

	return newUUID, nil
}
