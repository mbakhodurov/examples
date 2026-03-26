package kafka

import "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"

type UFORecordedDecoder interface {
	Decode(data []byte) (model.UFORecordedEvent, error)
}
