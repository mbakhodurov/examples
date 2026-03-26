package main

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/logger"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/config"

	kafkaConsumer "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/platform/pkg/kafka/consumer"
	kafkaConverter "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/converter/kafka/decoder"
	ufoConsumer "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/service/consumer/ufo_consumer"
)

const configPath = "../deploy/compose/ufo/.env"

func main() {
	// CONFIG
	if err := config.Load("../deploy/compose/ufo/.env"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// LOGGER
	if err := logger.Init("debug", false); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	// KAFKA CONFIG
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_6_0_0
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = false

	consumerGroup, err := sarama.NewConsumerGroup(
		config.AppConfig().Kafka.Brokers(),
		"ufo-consumer-group-1", // лучше новый group
		saramaConfig,
	)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	topics := []string{"ufo-recorded"}

	kConsumer := kafkaConsumer.NewConsumer(
		consumerGroup,
		topics,
		logger.Logger(),
	)

	decoder := kafkaConverter.NewUFORecordedDecoder()

	ufoConsumerService := ufoConsumer.NewService(kConsumer, decoder)

	go func() {
		if err := ufoConsumerService.RunConsumer(context.Background()); err != nil {
			log.Fatalf("consumer error: %v", err)
		}
	}()

	select {}
}
