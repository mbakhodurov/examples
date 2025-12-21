package ufo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/converter"
	repomodel "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, sighting model.SightingInfo) (string, error) {
	newUUID := uuid.NewString()
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[newUUID] = repomodel.Sighting{
		UUID:       newUUID,
		Info:       converter.SightingInfoToRepo(sighting),
		Created_at: time.Now(),
	}
	return r.data[newUUID].UUID, nil
}
