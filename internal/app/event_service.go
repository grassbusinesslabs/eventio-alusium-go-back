package app

import (
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"log"
	"time"
)

type EventService interface {
	Save(event domain.Event) (domain.Event, error)
	Update(event domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error)
	FindAll() ([]domain.Event, error)
	Delete(id uint64) error
	SubscribeToEvent(eventId, userId uint64) error
	GetUserSubscriptions(userId uint64) ([]domain.Event, error)
	FindEventsByDate(date time.Time) ([]domain.Event, error)
	FindEventsGroupByDate() (map[string][]domain.Event, error)
	FindList(filters database.UrlFilters) ([]domain.Event, error)
}

type eventService struct {
	subscriptionRepo database.SubscriptionRepository
	eventRepo        database.EventRepository
}

func NewEventService(ev database.EventRepository, sb database.SubscriptionRepository) EventService {
	return eventService{
		subscriptionRepo: sb,
		eventRepo:        ev,
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
func (s eventService) SubscribeToEvent(eventId, userId uint64) error {
	// Проверяем, существует ли событие
	_, err := s.eventRepo.Find(eventId)
	if err != nil {
		return err
	}

	// Добавляем подписку
	return s.subscriptionRepo.Subscribe(eventId, userId)
}

func (s eventService) GetUserSubscriptions(userId uint64) ([]domain.Event, error) {
	return s.subscriptionRepo.FindUserSubscriptions(userId)
}
func (s eventService) FindEventsByDate(date time.Time) ([]domain.Event, error) {
	_, err := s.eventRepo.FindEventsByDate(date)
	if err != nil {
		return nil, err
	}
	return s.eventRepo.FindEventsByDate(date)
}
func (s eventService) FindEventsGroupByDate() (map[string][]domain.Event, error) {
	return s.eventRepo.FindEventsGroupByDate()
}
func (s eventService) FindList(filters database.UrlFilters) ([]domain.Event, error) {
	events, err := s.eventRepo.FindList(filters)
	if err != nil {
		return nil, err
	}

	return events, nil
}
