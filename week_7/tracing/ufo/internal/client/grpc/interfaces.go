package grpc

import (
	"context"

	"github.com/mbakhodurov/examples/week_7/tracing/ufo/internal/model"
)

// AnalysisClient интерфейс для взаимодействия с Analysis сервисом
type AnalysisClient interface {
	AnalyzeSighting(ctx context.Context, uuid string) (model.AnalysisResult, error)
}
