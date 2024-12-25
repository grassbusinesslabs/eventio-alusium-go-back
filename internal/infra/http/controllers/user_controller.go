package controllers

import (
	"fmt"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/filesystem"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"io"
	"log"
	"net/http"
	"strconv"
)

type UserController struct {
	userService  app.UserService
	authService  app.AuthService
	imageService filesystem.ImageStorageService
}

func NewUserController(us app.UserService, as app.AuthService, imageService filesystem.ImageStorageService) UserController {
	return UserController{
		userService:  us,
		authService:  as,
		imageService: imageService,
	}
}

func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		Success(w, resources.UserDto{}.DomainToDto(user))
	}
}

func (c UserController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.UpdateUserRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}

		u := r.Context().Value(UserKey).(domain.User)
		u.FirstName = user.FirstName
		u.SecondName = user.SecondName
		u.Email = user.Email
		user, err = c.userService.Update(u)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		var userDto resources.UserDto
		Success(w, userDto.DomainToDto(user))
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)

		err := c.userService.Delete(u.Id)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
func (c UserController) SaveImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(UserKey).(domain.User)
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

		filename := fmt.Sprintf("user_%s_%s", strconv.FormatUint(u.Id, 10), header.Filename)

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

		u.Image = filename
		updatedUser, err := c.userService.Update(u)
		if err != nil {
			log.Printf("UserController -> SaveImage -> Update: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, map[string]string{"message": "File saved successfully!", "path": updatedUser.Image})
	}
}

func (c UserController) GetImage() http.HandlerFunc {
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
func (c UserController) DeleteImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, ok := r.Context().Value(UserKey).(domain.User)

		if !ok {
			InternalServerError(w, fmt.Errorf("failed to cast event"))
			return
		}

		if u.Image != "" {
			err := c.imageService.DeleteImage(u.Image)
			if err != nil {
				log.Printf("Failed to delete old image: %s", err)
				http.Error(w, "Failed to delete old image", http.StatusInternalServerError)
				return
			}
			u.Image = ""
			_, err = c.userService.Update(u)
			if err != nil {
				log.Printf("Failed to update event after image deletion: %s", err)
				InternalServerError(w, err)
				return
			}
		}
		Ok(w)
	}
}
func (c UserController) UpdateImage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, ok := r.Context().Value(UserKey).(domain.User)
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

		newFilename := fmt.Sprintf("user_%s_%s", strconv.FormatUint(u.Id, 10), header.Filename)

		if u.Image != "" {
			err = c.imageService.DeleteImage(u.Image)
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

		u.Image = newFilename
		updatedEvent, err := c.userService.Update(u)
		if err != nil {
			log.Printf("UserController -> SaveImage -> Update: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, map[string]string{"message": "File saved successfully!", "path": updatedEvent.Image})
	}
}
