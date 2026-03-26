package ufo_consumer

import (
	"context"

	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	kafkaConverter "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/converter/kafka"
	def "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service"
	"go.uber.org/zap"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	ufoRecordedConsumer kafka.Consumer
	ufoRecordedDecoder  kafkaConverter.UFORecordedDecoder
}

func NewService(ufoRecordedConsumer kafka.Consumer, ufoRecordedDecoder kafkaConverter.UFORecordedDecoder) *service {
	return &service{
		ufoRecordedConsumer: ufoRecordedConsumer,
		ufoRecordedDecoder:  ufoRecordedDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order ufoRecordedConsumer service")

	if err := s.ufoRecordedConsumer.Consume(ctx, s.OrderHandler); err != nil {
		logger.Error(ctx, "Consume from ufo.recorded topic error", zap.Error(err))
		return err
	}
	return nil
}
