package ufo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mbakhodurov/examples/week_7/tracing/ufo/internal/model"
	repoConverter "github.com/mbakhodurov/examples/week_7/tracing/ufo/internal/repository/converter"
	repoModel "github.com/mbakhodurov/examples/week_7/tracing/ufo/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	var repoSighting repoModel.Sighting

	if err := r.collection.FindOne(
		ctx,
		bson.M{"uuid": uuid},
	).Decode(&repoSighting); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Sighting{}, model.ErrSightingNotFound
		}

		return model.Sighting{}, err
	}

	return repoConverter.SightingToModel(repoSighting), nil
}
