package database

import (
	"fmt"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
	"log"
)

type subscriptionRepository struct {
	db db.Session
}

type SubscriptionRepository interface {
	Subscribe(eventId, userId uint64) error
	FindUserSubscriptions(userId uint64) ([]domain.Event, error)
}

func NewSubscriptionRepository(db db.Session) SubscriptionRepository {
	return subscriptionRepository{db: db}
}

func (r subscriptionRepository) Subscribe(eventId, userId uint64) error {
	var event map[string]interface{}
	err := r.db.Collection("events").
		Find(db.Cond{"id": eventId, "deleted_date": nil}).One(&event)

	if err != nil || event == nil {
		log.Printf("SubscribeRepository -> Subscribe -> r.db.Collection(\"events\"): %s", err)
		return fmt.Errorf("event not found or is deleted")
	}

	_, er := r.db.Collection("subscriptions").Insert(map[string]interface{}{
		"event_id": eventId,
		"user_id":  userId,
	})
	return er
}

func (r subscriptionRepository) FindUserSubscriptions(userId uint64) ([]domain.Event, error) {
	var dbEvents []event

	// Выполняем запрос
	err := r.db.SQL().
		Select("e.id", "e.user_id", "e.title", "e.description", "e.status", "e.image", "e.location", "e.date", "e.lat", "e.lon").
		From("subscriptions AS s").
		Join("events AS e").On("s.event_id = e.id").
		Where("s.user_id = ? AND e.deleted_date IS NULL", userId).
		All(&dbEvents)

	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(dbEvents), nil
}
func (r subscriptionRepository) mapModelToDomain(m event) domain.Event {
	return domain.Event{
		Id:          m.Id,
		UserId:      m.UserId,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		Image:       m.Image,
		Location:    m.Location,
		Lat:         m.Lat,
		Lon:         m.Lon,
		Date:        m.Date,
	}
}
func (r subscriptionRepository) mapModelToDomainCollection(evn []event) []domain.Event {

	var events []domain.Event
	for _, ev := range evn {

		events = append(events, r.mapModelToDomain(ev))
	}
	return events
}
