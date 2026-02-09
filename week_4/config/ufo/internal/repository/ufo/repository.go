package ufo

import (
	def "github.com/mbakhodurov/examples/week_4/config/ufo/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ def.UFORepository = (*repository)(nil)

const (
	databaseName   = "ufo_db"
	collectionName = "sightings"
)

type repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client) *repository {
	return &repository{
		client:     client,
		collection: client.Database(databaseName).Collection(collectionName),
	}
}
