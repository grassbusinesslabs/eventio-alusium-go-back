package app

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"log"
)

type EventService interface {
	Save(event domain.Event) (domain.Event, error)
	Update(event domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error)
	FindAll() ([]domain.Event, error)
	Delete(id uint64) error
}

type eventService struct {
	eventRepo database.EventRepository
}

func NewEventService(ev database.EventRepository) EventService {
	return eventService{
		eventRepo: ev,
	}
}

func (s eventService) Save(e domain.Event) (domain.Event, error) {
	evn, err := s.eventRepo.Save(e)
	if err != nil {
		log.Printf("EventService -> Save -> s.eventRepo.Save: %s", err)
		return domain.Event{}, err
	}
	return evn, nil
}
func (s eventService) Update(event domain.Event) (domain.Event, error) {
	event, err := s.eventRepo.Update(event)
	if err != nil {
		log.Printf("Event service -> Update -> s.eventRepo.Update(event): %s", err)
		return domain.Event{}, err
	}

	return event, nil
}
func (s eventService) Find(id uint64) (interface{}, error) {
	event, err := s.eventRepo.Find(id)
	if err != nil {
		log.Printf("Event service -> Find -> s.eventRepo.Find(id): %s", err)
		return domain.Event{}, err
	}

	return event, nil
}
func (s eventService) FindAll() ([]domain.Event, error) {
	events, err := s.eventRepo.FindAll()
	if err != nil {
		log.Printf("Event service -> FindAll -> events, s.eventRepo.FindAll(): %s", err)
		return nil, err
	}
	return events, nil
}
func (s eventService) Delete(id uint64) error {
	err := s.eventRepo.Delete(id)
	if err != nil {
		log.Printf("EventService -> Delete -> s.eventRepo.Delete(id): %s", err)
		return err
	}

	return nil
}
