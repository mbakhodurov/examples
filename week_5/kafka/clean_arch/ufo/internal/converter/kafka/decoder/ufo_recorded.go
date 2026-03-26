package decoder

import (
	"fmt"
	"time"

	events_v1 "github.com/mbakhodurov/examples/week_5/kafka/clean_arch/shared/pkg/proto/events/v1"
	"github.com/mbakhodurov/examples/week_5/kafka/clean_arch/ufo/internal/model"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewUFORecordedDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.UFORecordedEvent, error) {
	var pb events_v1.UFORecorded
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.UFORecordedEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var observedAt *time.Time
	if pb.ObservedAt != nil {
		observedAt = lo.ToPtr(pb.ObservedAt.AsTime())
	}

	return model.UFORecordedEvent{
		UUID:        pb.Uuid,
		ObservedAt:  observedAt,
		Location:    pb.Location,
		Description: pb.Description,
	}, nil
}
