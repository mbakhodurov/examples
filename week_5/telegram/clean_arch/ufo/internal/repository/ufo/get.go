package ufo

import (
	"context"
	"errors"

	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/model"
	repoconverter "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository/converter"
	repomodels "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
