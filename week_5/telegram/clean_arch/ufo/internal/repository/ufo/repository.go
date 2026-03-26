package ufo

import (
	def "github.com/mbakhodurov/examples/week_5/telegram/clean_arch/ufo/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "sightings"
)

var _ def.UFORepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection(collectionName),
	}
}
