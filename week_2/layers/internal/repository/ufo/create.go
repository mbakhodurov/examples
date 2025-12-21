package ufo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/examples/week_2/layers/internal/model"
	repoConverter "github.com/mbakhodurov/examples/week_2/layers/internal/repository/converter"
	repoModel "github.com/mbakhodurov/examples/week_2/layers/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	newUUID := uuid.NewString()
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[newUUID] = repoModel.Sighting{
		Uuid:       newUUID,
		Info:       repoConverter.SightingInfoToRepoModel(info),
		Created_at: time.Now(),
	}
	return newUUID, nil
}
