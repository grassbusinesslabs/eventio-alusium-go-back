package resources

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"time"
)

type EventDto struct {
	Id          uint64             `db:"id,omitempty"`
	UserId      uint64             `db:"user_id,omitempty"`
	Title       string             `db:"title"`
	Description string             `db:"description"`
	Status      domain.EventStatus `db:"status"`
	Date        time.Time          `db:"date"`
	Image       string             `db:"image"`
	Location    string             `db:"location"`
	Lat         float64            `db:"lat"`
	Lon         float64            `db:"long"`
}
type EventsDto struct {
	Events []EventDto `json:"events"`
}

func (d EventsDto) DomainToDto(ev []domain.Event) EventsDto {
	events := make([]EventDto, len(ev))
	for i, e := range ev {
		events[i] = EventDto{}.DomainToDto(e)
	}

	return EventsDto{
		Events: events,
	}
}

func (d EventDto) DomainToDto(event domain.Event) EventDto {
	return EventDto{
		Id:          event.Id,
		UserId:      event.UserId,
		Title:       event.Title,
		Description: event.Description,
		Status:      event.Status,
		Image:       event.Image,
		Location:    event.Location,
		Lat:         event.Lat,
		Lon:         event.Lon,
		Date:        event.Date,
	}
}
