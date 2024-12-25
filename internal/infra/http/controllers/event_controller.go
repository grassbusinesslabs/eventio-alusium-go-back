package controllers

import (
	"fmt"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/filesystem"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type EventController struct {
	eventService app.EventService
	imageService filesystem.ImageStorageService
}

func NewEventController(eventService app.EventService, imageService filesystem.ImageStorageService) EventController {
	return EventController{
		eventService: eventService,
		imageService: imageService,
	}
}

func (c EventController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := requests.Bind(r, requests.CreateEventRequest{}, domain.Event{})
		if err != nil {
			log.Printf("EventController -> Save -> requests.Bind: %s", err)
			BadRequest(w, err)
			return
		}
		user := r.Context().Value(UserKey).(domain.User)

		event.UserId = user.Id
		event.Status = domain.NewEventStatus

		event, err = c.eventService.Save(event)
		if err != nil {
			log.Printf("EventController -> Save -> c.eventService.Save(event): %s", err)
			InternalServerError(w, err)
			return
		}
		var eventDto resources.EventDto
		eventDto = eventDto.DomainToDto(event)
		Created(w, eventDto)
	}
}
func (c EventController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqevent, err := requests.Bind(r, requests.UpdateEventRequest{}, domain.Event{})
		if err != nil {
			log.Printf("EventController -> Update -> requests.Bind:%s", err)
			BadRequest(w, err)
			return
		}
		ev, ok := r.Context().Value(EventKey).(domain.Event)
		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}
		ev.Title = reqevent.Title
		ev.Description = reqevent.Description
		reqevent, err = c.eventService.Update(ev)

		if err != nil {
			log.Printf("EventController -> Update -> c.eventService.Update(ev): %s", err)
			InternalServerError(w, err)
			return
		}

		var eventDto resources.EventDto
		Success(w, eventDto.DomainToDto(reqevent))
	}
}
func (c EventController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ev, ok := r.Context().Value(EventKey).(domain.Event)
		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}

		er := c.eventService.Delete(ev.Id)
		if er != nil {
			log.Printf("EventController -> Delete -> c.eventService.Delete(eventId): %s", er)
			InternalServerError(w, er)
			return
		}

		Ok(w)
	}
}
func (c EventController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		events, err := c.eventService.FindAll()

		if err != nil {
			log.Printf("EventController -> FindAll -> evns, err := c.eventService.FindAll(): %s", err)
			InternalServerError(w, err)
			return
		}
		var eventsDto resources.EventsDto
		Success(w, eventsDto.DomainToDto(events))
		Ok(w)
	}
}
func (c EventController) Subscribe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ev, ok := r.Context().Value(EventKey).(domain.Event)
		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		if err := c.eventService.SubscribeToEvent(ev.Id, user.Id); err != nil {
			InternalServerError(w, err)
			return
		}

		Success(w, "Subscribed successfully")
	}
}
func (c EventController) GetUserSubscriptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		subscriptions, err := c.eventService.GetUserSubscriptions(user.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		var eventsDto resources.EventsDto
		Success(w, eventsDto.DomainToDto(subscriptions))

	}
}

func (c EventController) FindEventsByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateParam := r.URL.Query().Get("date")
		if dateParam == "" {
			BadRequest(w, fmt.Errorf("missing date parameter"))
			return
		}

		timestamp, err := strconv.ParseInt(dateParam, 10, 64)
		if err != nil {
			BadRequest(w, fmt.Errorf("invalid date format, expected Unix timestamp"))
			return
		}

		date := time.Unix(timestamp, 0)

		events, err := c.eventService.FindEventsByDate(date)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		var eventsDto resources.EventsDto
		Success(w, eventsDto.DomainToDto(events))

	}
}

func (c EventController) FindEventsGroupByDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupedEvents, err := c.eventService.FindEventsGroupByDate()
		if err != nil {
			InternalServerError(w, err)
			return
		}

		filteredGroupedEvents := make(map[string][]resources.EventDto)

		for date, events := range groupedEvents {
			var filteredEvents []resources.EventDto
			for _, event := range events {
				filteredEvents = append(filteredEvents, resources.EventDto{}.DomainToDto(event))
			}
			filteredGroupedEvents[date] = filteredEvents
		}

		// Отправляем результат
		Success(w, filteredGroupedEvents)
	}
}
func (c EventController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		search := r.URL.Query().Get("search")
		location := r.URL.Query().Get("location")
		dateParam := r.URL.Query().Get("date")
		city := r.URL.Query().Get("city")
		var date *time.Time
		if dateParam != "" {
			timestamp, err := strconv.ParseInt(dateParam, 10, 64)
			if err != nil {
				http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
			parsedDate := time.Unix(timestamp, 0)
			date = &parsedDate
		}

		filters := database.UrlFilters{
			Search:   search,
			Location: location,
			Date:     date,
			City:     city,
		}

		events, err := c.eventService.FindList(filters)
		if err != nil {
			http.Error(w, "Error fetching events", http.StatusInternalServerError)
			return
		}

		var eventsDto resources.EventsDto
		Success(w, eventsDto.DomainToDto(events))

	}
}

func (c EventController) SaveImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ev, ok := r.Context().Value(EventKey).(domain.Event)
		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to get the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		filename := fmt.Sprintf("event_%s_%s", strconv.FormatUint(ev.Id, 10), header.Filename)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read the file", http.StatusInternalServerError)
			return
		}

		err = c.imageService.SaveImage(filename, content)
		if err != nil {
			http.Error(w, "Failed to save the image", http.StatusInternalServerError)
			return
		}

		ev.Image = filename
		updatedEvent, err := c.eventService.Update(ev)
		if err != nil {
			log.Printf("EventController -> UploadImage -> Update: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, map[string]string{"message": "File saved successfully!", "path": updatedEvent.Image})
	}
}
func (c EventController) GetImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imgPath := r.URL.Query().Get("path")
		content, err := c.imageService.GetImageContent(imgPath)
		if err != nil {
			http.Error(w, "Failed to get the image", http.StatusInternalServerError)
			return
		}

		Success(w, content)
	}
}
func (c EventController) DeleteImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ev, ok := r.Context().Value(EventKey).(domain.Event)

		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}

		if ev.Image != "" {
			err := c.imageService.DeleteImage(ev.Image)
			if err != nil {
				log.Printf("Failed to delete old image: %s", err)
				http.Error(w, "Failed to delete old image", http.StatusInternalServerError)
				return
			}
			ev.Image = ""
			_, err = c.eventService.Update(ev)
			if err != nil {
				log.Printf("Failed to update event after image deletion: %s", err)
				InternalServerError(w, err)
				return
			}
		}
		Ok(w)
	}
}

func (c EventController) UpdateImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ev, ok := r.Context().Value(EventKey).(domain.Event)
		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}

		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to get the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		newFilename := fmt.Sprintf("event_%s_%s", strconv.FormatUint(ev.Id, 10), header.Filename)

		if ev.Image != "" {
			err = c.imageService.DeleteImage(ev.Image)
			if err != nil {
				log.Printf("Failed to delete old image: %s", err)
				http.Error(w, "Failed to delete old image", http.StatusInternalServerError)
				return
			}
		}

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read the file", http.StatusInternalServerError)
			return
		}

		err = c.imageService.SaveImage(newFilename, content)
		if err != nil {
			http.Error(w, "Failed to save the image", http.StatusInternalServerError)
			return
		}

		ev.Image = newFilename
		updatedEvent, err := c.eventService.Update(ev)
		if err != nil {
			log.Printf("EventController -> SaveImage -> Update: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, map[string]string{"message": "File saved successfully!", "path": updatedEvent.Image})
	}
}
