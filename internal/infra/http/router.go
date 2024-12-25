package http

import (
	"encoding/json"
	"fmt"
	"github.com/BohdanBoriak/boilerplate-go-back/config"
	"github.com/BohdanBoriak/boilerplate-go-back/config/container"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "capacitor://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController)
				EventRouter(apiRouter, cont.EventController, cont.PathMw)
				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, ac controllers.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {

		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
		apiRouter.Post(
			"/saveImage",
			uc.SaveImage(),
		)
		apiRouter.Get(
			"/getImage",
			uc.GetImage(),
		)
		apiRouter.Delete(
			"/deleteImage",
			uc.DeleteImage(),
		)
		apiRouter.Post(
			"/updateImage",
			uc.UpdateImage(),
		)
	})
}

func EventRouter(r chi.Router, ev controllers.EventController, pathMw func(http.Handler) http.Handler) {
	r.Route("/events", func(apiRouter chi.Router) {

		apiRouter.Post(
			"/",
			ev.Save(),
		)
		apiRouter.With(pathMw).Put(
			"/update/{eventId}",
			ev.Update(),
		)
		apiRouter.With(pathMw).Delete(
			"/delete/{eventId}",
			ev.Delete(),
		)
		apiRouter.Get(
			"/findAll",
			ev.FindAll(),
		)
		apiRouter.With(pathMw).Post(
			"/subscribe/{eventId}", ev.Subscribe(),
		)
		apiRouter.Get(
			"/subscriptions", ev.GetUserSubscriptions(),
		)
		apiRouter.Get(
			"/findByDate", ev.FindEventsByDate(),
		)
		apiRouter.Get(
			"/groupedByDate",
			ev.FindEventsGroupByDate(),
		)
		apiRouter.Get(
			"/findList",
			ev.FindList(),
		)
		apiRouter.With(pathMw).Post(
			"/{eventId}/uploadImage",
			ev.SaveImage(),
		)
		apiRouter.Get(
			"/image",
			ev.GetImage(),
		)
		apiRouter.With(pathMw).Delete(
			"/{eventId}/deleteImage",
			ev.DeleteImage(),
		)
		apiRouter.With(pathMw).Post(
			"/{eventId}/updateImage",
			ev.UpdateImage(),
		)
	})
}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
