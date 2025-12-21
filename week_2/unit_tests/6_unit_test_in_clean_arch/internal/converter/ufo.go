package converter

import (
	"time"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/model"
	ufo_v1 "github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func SightingModelToProto(sighting model.Sighting) *ufo_v1.Sighting {
	var updatedAt *timestamppb.Timestamp
	if sighting.Updated_at != nil {
		updatedAt = timestamppb.New(*sighting.Updated_at)
	}

	var deletedAt *timestamppb.Timestamp
	if sighting.Deleted_at != nil {
		deletedAt = timestamppb.New(*sighting.Deleted_at)
	}

	return &ufo_v1.Sighting{
		Uuid:      sighting.UUID,
		Info:      SightingInfoModelToProto(sighting.Info),
		CreatedAt: timestamppb.New(sighting.Created_at),
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func SightingInfoModelToProto(info model.SightingInfo) *ufo_v1.SightingInfo {
	var observedAt *timestamppb.Timestamp
	if info.Observed_at != nil {
		observedAt = timestamppb.New(*info.Observed_at)
	}

	var color *wrapperspb.StringValue
	if info.Color != nil {
		color = wrapperspb.String(*info.Color)
	}

	var sound *wrapperspb.BoolValue
	if info.Sound != nil {
		sound = wrapperspb.Bool(*info.Sound)
	}

	var durationSeconds *wrapperspb.Int32Value
	if info.Duration_seconds != nil {
		durationSeconds = wrapperspb.Int32(*info.Duration_seconds)
	}

	return &ufo_v1.SightingInfo{
		ObservedAt:      observedAt,
		Location:        info.Location,
		Description:     info.Description,
		Color:           color,
		Sound:           sound,
		DurationSeconds: durationSeconds,
	}
}

func SightingInfoProtoToModel(info *ufo_v1.SightingInfo) model.SightingInfo {
	var observed_at *time.Time
	if info.ObservedAt != nil {
		observed_at = lo.ToPtr(info.ObservedAt.AsTime())
	}

	var color *string
	if info.Color != nil {
		color = lo.ToPtr(info.Color.Value)
	}

	var sound *bool
	if info.Sound != nil {
		tmp := info.Sound.Value
		sound = &tmp
	}

	var durationSeconds *int32
	if info.DurationSeconds != nil {
		tmp := info.DurationSeconds.Value
		durationSeconds = &tmp
	}

	return model.SightingInfo{
		Observed_at:      observed_at,
		Location:         info.Location,
		Description:      info.Description,
		Color:            color,
		Sound:            sound,
		Duration_seconds: durationSeconds,
	}
}

func SightingUpdateInfoToModel(info *ufo_v1.SightingUpdateInfo) model.SightingUpdateInfo {
	var observedAt *time.Time
	if info.ObservedAt != nil {
		tmp := info.ObservedAt.AsTime()
		observedAt = &tmp
	}

	var location *string
	if info.Location != nil {
		tmp := info.Location.Value
		location = &tmp
	}

	var description *string
	if info.Description != nil {
		tmp := info.Description.Value
		description = &tmp
	}

	var color *string
	if info.Color != nil {
		tmp := info.Color.Value
		color = &tmp
	}

	var sound *bool
	if info.Sound != nil {
		tmp := info.Sound.Value
		sound = &tmp
	}

	var durationSeconds *int32
	if info.DurationSeconds != nil {
		tmp := info.DurationSeconds.Value
		durationSeconds = &tmp
	}
	return model.SightingUpdateInfo{
		Observed_at:      observedAt,
		Location:         location,
		Description:      description,
		Color:            color,
		Sound:            sound,
		Duration_seconds: durationSeconds,
	}
}
