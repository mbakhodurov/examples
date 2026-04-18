package model

import "time"

type SightingInfo struct {
	ObservedAt      *time.Time
	Location        string
	Description     string
	Color           *string
	Sound           *bool
	DurationSeconds *int32
}

type SightingUpdateInfo struct {
	ObservedAt      *time.Time
	Location        *string
	Description     *string
	Color           *string
	Sound           *bool
	DurationSeconds *int32
}

type Sighting struct {
	Uuid      string
	Info      SightingInfo
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type SightingRedisView struct {
	Uuid            string  `redis:"uuid"`
	ObservedAtNs    *int64  `redis:"observed_at,omitempty"`
	Location        string  `redis:"location"`
	Description     string  `redis:"description"`
	Color           *string `redis:"color,omitempty"`
	Sound           *bool   `redis:"sound,omitempty"`
	DurationSeconds *int32  `redis:"duration_seconds,omitempty"`
	CreatedAtNs     int64   `redis:"created_at"`
	UpdatedAtNs     *int64  `redis:"updated_at,omitempty"`
	DeletedAtNs     *int64  `redis:"deleted_at,omitempty"`
}
