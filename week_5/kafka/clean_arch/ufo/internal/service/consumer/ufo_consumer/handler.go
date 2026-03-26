package ufo_consumer

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.ufoRecordedDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode UFORecorded", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.UUID),
		zap.String("location", event.Location),
		zap.String("description", event.Description),
	)

	return nil
}
