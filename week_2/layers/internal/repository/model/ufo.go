package model

import "time"

type SightingInfo struct {
	Observed_at      *time.Time
	Location         string
	Description      string
	Color            *string
	Sound            *bool
	Duration_seconds *int32
}

type Sighting struct {
	Uuid       string
	Info       SightingInfo
	Created_at time.Time
	Updated_at *time.Time
	Deleted_at *time.Time
}

type SightingUpdateInfo struct {
	Observed_at      *time.Time
	Location         *string
	Description      *string
	Color            *string
	Sound            *bool
	Duration_seconds *int32
}
