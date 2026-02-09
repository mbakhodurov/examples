package ufo

import (
	"context"
	"errors"

	repoconverter "github.com/mbakhodurov/examples/week_4/di/ufo/internal/repository/converter"
	repomodels "github.com/mbakhodurov/examples/week_4/di/ufo/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	var repoModel repomodels.Sighting

	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(repoModel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Sighting{}, model.ErrSightingNotFound
		}
		return model.Sighting{}, err
	}
	return repoconverter.SightingToModel(repoModel), nil
}
