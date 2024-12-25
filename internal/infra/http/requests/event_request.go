package requests

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"time"
)

type CreateEventRequest struct {
	Title       string  `json:"title" validate:"required,max=80"`
	Description string  `json:"description"  validate:"required,max=200"`
	Image       string  `json:"image"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
	City        string  `json:"city"`
	Location    string  `json:"location"  validate:"required,max=200"`
	Date        int64   `json:"date"`
}
type UpdateEventRequest struct {
	Title       string  `json:"title" validate:"required,max=40"`
	Description string  `json:"description"  validate:"required,max=200"`
	Image       string  `json:"image" validate:"required"`
	Lat         float64 `json:"lat" validate:"required"`
	Lon         float64 `json:"lon" validate:"required"`
	City        string  `json:"city"`
	Location    string  `json:"location"  validate:"required,max=200"`
	Date        int64   `json:"date"`
}

func (r CreateEventRequest) ToDomainModel() (interface{}, error) {
	return domain.Event{
		Title:       r.Title,
		Description: r.Description,
		Image:       r.Image,
		City:        r.City,
		Location:    r.Location,
		Lat:         r.Lat,
		Lon:         r.Lon,
		Date:        time.Unix(r.Date, 0),
	}, nil
}

func (r UpdateEventRequest) ToDomainModel() (interface{}, error) {
	return domain.Event{
		Title:       r.Title,
		Description: r.Description,
		Image:       r.Image,
		City:        r.City,
		Location:    r.Location,
		Lat:         r.Lat,
		Lon:         r.Lon,
		Date:        time.Unix(r.Date, 0),
	}, nil
}
