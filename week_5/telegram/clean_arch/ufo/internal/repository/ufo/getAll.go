package ufo

import (
	"context"
	"errors"
	"log"

	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/model"
	"github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository/converter"
	repoModel "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetAll(ctx context.Context) ([]model.Sighting, error) {

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrSightingNotFound
		}
		return nil, err
	}

	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var sightings []model.Sighting

	for cursor.Next(ctx) {
		var repoSight repoModel.Sighting

		if err := cursor.Decode(&repoSight); err != nil {
			return nil, err
		}
		sightings = append(sightings, converter.SightingToModel(repoSight))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(sightings) == 0 {
		return nil, model.ErrSightingNotFound
	}

	return sightings, nil
}
