package container

import (
	"github.com/BohdanBoriak/boilerplate-go-back/config"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/controllers"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/middlewares"
	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"net/http"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw func(http.Handler) http.Handler
	PathMw func(http.Handler) http.Handler
}

type Services struct {
	app.AuthService
	app.UserService
}

type Controllers struct {
	AuthController  controllers.AuthController
	UserController  controllers.UserController
	EventController controllers.EventController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	sessionRepository := database.NewSessRepository(sess)
	userRepository := database.NewUserRepository(sess)
	eventRepository := database.NewEventRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userRepository, tknAuth, conf.JwtTTL)
	eventService := app.NewEventService(eventRepository)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService, authService)
	eventController := controllers.NewEventController(eventService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)
	pathObjMiddleware := middlewares.PathObject("eventId", controllers.EventKey, eventService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
			PathMw: pathObjMiddleware,
		},
		Services: Services{
			authService,
			userService,
		},
		Controllers: Controllers{
			authController,
			userController,
			eventController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
