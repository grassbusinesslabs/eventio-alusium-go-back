package domain

import (
	"time"
)

type Event struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	Status      EventStatus
	Image       string
	Location    string
	Date        time.Time
	Lat         float64
	Lon         float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type EventStatus string

const (
	NewEventStatus  EventStatus = "NEW"
	DoneEventStatus EventStatus = "DONE"
)

func (e Event) GetEventId() uint64 {
	return e.Id
}
