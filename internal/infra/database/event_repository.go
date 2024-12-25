package database

import (
	"github.com/upper/db/v4"
	"log"
	"strings"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

const EventTableName = "events"

type event struct {
	Id          uint64             `db:"id,omitempty"`
	UserId      uint64             `db:"user_id,omitempty"`
	Title       string             `db:"title"`
	Description string             `db:"description"`
	Status      domain.EventStatus `db:"status"`
	Image       string             `db:"image"`
	City        string             `db:"city"`
	Location    string             `db:"location"`
	Lat         float64            `db:"lat"`
	Lon         float64            `db:"lon"`
	Date        time.Time          `db:"date"`
	CreatedDate time.Time          `db:"created_date,omitempty"`
	UpdatedDate time.Time          `db:"updated_date,omitempty"`
	DeletedDate *time.Time         `db:"deleted_date,omitempty"`
}
type EventRepository interface {
	Save(event domain.Event) (domain.Event, error)
	Update(event domain.Event) (domain.Event, error)
	Find(id uint64) (interface{}, error)
	Delete(id uint64) error
	FindAll() ([]domain.Event, error)
	FindEventsByDate(date time.Time) ([]domain.Event, error)
	FindEventsGroupByDate() (map[string][]domain.Event, error)
	FindList(filters UrlFilters) ([]domain.Event, error)
}

type eventRepository struct {
	coll db.Collection
	sess db.Session
}

type UrlFilters struct {
	Search   string
	Date     *time.Time
	Location string
	City     string
}

func NewEventRepository(dbSession db.Session) eventRepository {
	return eventRepository{
		coll: dbSession.Collection(EventTableName),
		sess: dbSession,
	}
}

func (r eventRepository) Save(event domain.Event) (domain.Event, error) {
	evn := r.mapDomainToModel(event)
	evn.CreatedDate, evn.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&evn)
	if err != nil {
		log.Printf("EventRepository -> Save -> r.coll.InsertReturning(&evn) %s", err)
		return domain.Event{}, err
	}
	return r.mapModelToDomain(evn), nil
}

func (r eventRepository) Update(event domain.Event) (domain.Event, error) {
	e := r.mapDomainToModel(event)
	e.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": e.Id, "deleted_date": nil}).Update(&e)
	if err != nil {
		log.Printf("EventRepository -> Update -> r.coll.Find(db.Cond{\"id\": e.Id, \"deleted_date\": nil}).Update(&e): %s", err)
		return domain.Event{}, err
	}
	return r.mapModelToDomain(e), nil
}

func (r eventRepository) Find(id uint64) (interface{}, error) {
	var evn event
	err := r.coll.Find(db.Cond{"id": id}).One(&evn)
	if err != nil {
		log.Printf("EventRepository -> Find -> r.coll.Find(db.Cond{\"id\": id}).One(&evn) %s", err)
		return nil, err
	}

	return r.mapModelToDomain(evn), nil
}
func (r eventRepository) FindAll() ([]domain.Event, error) {
	var events []event
	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&events)
	if err != nil {
		log.Printf("EventRepository -> FindAll -> r.coll.Find(db.Cond{\"deleted_date\": nil}).All(&events) %s", err)
		return []domain.Event{}, err
	}
	return r.mapModelToDomainCollection(events), nil
}
func (r eventRepository) FindEventsByDate(date time.Time) ([]domain.Event, error) {
	var events []event
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	err := r.coll.Find(db.Cond{
		"date >=":      startOfDay,
		"date <=":      endOfDay,
		"deleted_date": nil,
	}).All(&events)

	if err != nil {
		log.Printf("EventRepository -> FindEventsByDate -> r.coll.Find(db.Cond{\n\t\t\"date >=\":startOfDay,\n\t\t\"date <=\":endOfDay,\n\t\t\"deleted_date\": nil,\n\t}) %s", err)
		return nil, err
	}

	return r.mapModelToDomainCollection(events), nil
}
func (r eventRepository) FindEventsGroupByDate() (map[string][]domain.Event, error) {
	groupedEvents := make(map[string][]domain.Event)

	var events []event
	err := r.coll.Find(db.Cond{"deleted_date": nil}).All(&events)
	if err != nil {
		log.Printf("EventRepository -> FindEventsGroupByDate -> r.coll.Find(db.Cond{\"deleted_date\": nil}).All(&events) %s", err)
		return nil, err
	}

	for _, e := range events {
		dateKey := e.Date.Format("2006-01-02")
		groupedEvents[dateKey] = append(groupedEvents[dateKey], r.mapModelToDomain(e))
	}

	return groupedEvents, nil
}
func (r eventRepository) FindList(filters UrlFilters) ([]domain.Event, error) {
	query := r.coll.Find(db.Cond{"deleted_date": nil})

	if filters.City != "" {
		city := "%" + strings.ToLower(filters.City) + "%"
		query = query.And(db.Raw(`LOWER(city) LIKE ?`, city))
	}

	if filters.Search != "" {
		search := "%" + strings.ToLower(filters.Search) + "%"
		query = query.And(db.Raw(`(LOWER(title) LIKE ? OR LOWER(description) LIKE ?)`, search, search))
	}

	if filters.Date != nil {
		startOfDay := filters.Date.Truncate(24 * time.Hour)
		endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)
		query = query.And(db.Cond{"date >=": startOfDay, "date <=": endOfDay})
	}

	if filters.Location != "" {
		location := "%" + strings.ToLower(filters.Location) + "%"
		query = query.And(db.Raw(`LOWER(location) LIKE ?`, location))
	}

	var events []event
	err := query.OrderBy("-date").All(&events)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(events), nil

}

func (r eventRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r eventRepository) mapDomainToModel(d domain.Event) event {
	return event{
		Id:          d.Id,
		UserId:      d.UserId,
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
		Image:       d.Image,
		City:        d.City,
		Location:    d.Location,
		Lat:         d.Lat,
		Lon:         d.Lon,
		Date:        d.Date,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r eventRepository) mapModelToDomain(m event) domain.Event {
	return domain.Event{
		Id:          m.Id,
		UserId:      m.UserId,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		Image:       m.Image,
		City:        m.City,
		Location:    m.Location,
		Lat:         m.Lat,
		Lon:         m.Lon,
		Date:        m.Date,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
func (r eventRepository) mapModelToDomainCollection(evn []event) []domain.Event {

	var events []domain.Event
	for _, ev := range evn {

		events = append(events, r.mapModelToDomain(ev))
	}
	return events
}
