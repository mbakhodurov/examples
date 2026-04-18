package ufo

import (
	def "github.com/mbakhodurov/examples/week_6/redis/clean_arch/ufo/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ def.UFORepository = (*repository)(nil)

const (
	collectionName = "sightings"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection(collectionName),
	}
}
