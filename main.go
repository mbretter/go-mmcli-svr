package main

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mbretter/go-mmcli-svr/api"
	"github.com/mbretter/go-mmcli-svr/backend/mmcli"
	_ "github.com/mbretter/go-mmcli-svr/docs"
	"github.com/mbretter/go-mmcli-svr/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
)

const (
	AppModeTest string = "tst"
	AppModeProd string = "prod"
	AppModeDev  string = "dev"
)

// @title		mmcli server
// @version		1.0
// @description	a http server in front of mmcli
//
// @license.name	bsd
// @host	127.0.0.1:8743
//
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load(".env")
	_ = godotenv.Overload(".env.local")

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Debug("started")

	backend := mmcli.Provide()

	commandLine := NewCommandLine(log, backend)
	err := commandLine.Parse()
	if err != nil {
		panic(err)
	}

	commandLine.Activate()

	r := chi.NewRouter()
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.NoCache)

	r.Use(middleware.HttpLoggerMiddleware(log, chiMiddleware.RequestIDKey))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization", "Accept-Language"},
		ExposedHeaders:   []string{"Content-Length", "X-Message"},
		AllowCredentials: true,
		MaxAge:           7200,
		Debug:            false,
	}))

	r.Use(middleware.HttpModemMiddleware(log))

	r.With(middleware.LogRoute).Route("/", func(r chi.Router) {
		handlers := api.Provide(backend)
		registerModemRoutes(r, handlers)
		registerLocationRoutes(r, handlers)
		registerSmsRoutes(r, handlers)

		utilsApi := api.ProvideUtilsApi()
		r.Get("/ping", utilsApi.Ping)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		))
	})

	log.Info("Listening", "address", commandLine.Listen)
	//goland:noinspection ALL
	http.ListenAndServe(commandLine.Listen, r)
}

type modemHandlersInterface interface {
	ModemList(w http.ResponseWriter, r *http.Request)
	ModemDetail(w http.ResponseWriter, r *http.Request)
}

func registerModemRoutes(r chi.Router, handlers modemHandlersInterface) {
	r.Get("/modem/", handlers.ModemList)
	r.Get("/modem/{id:[a-zA-Z0-9%/]+}", handlers.ModemDetail)
}

type locationHandlersInterface interface {
	LocationGet(w http.ResponseWriter, r *http.Request)
}

func registerLocationRoutes(r chi.Router, handlers locationHandlersInterface) {
	r.Get("/location", handlers.LocationGet)
}

type smsHandlersInterface interface {
	SmsGet(w http.ResponseWriter, r *http.Request)
	SmsSend(w http.ResponseWriter, r *http.Request)
	SmsDelete(w http.ResponseWriter, r *http.Request)
}

func registerSmsRoutes(r chi.Router, handlers smsHandlersInterface) {
	r.Get("/sms/", handlers.SmsGet)
	r.Get("/sms/{id:[a-zA-Z0-9%/]+}", handlers.SmsGet)
	r.Post("/sms", handlers.SmsSend)
	r.Delete("/sms/{id:[a-zA-Z0-9%/]+}", handlers.SmsDelete)
}
