package main

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go-mmcli-svr/api"
	be "go-mmcli-svr/backend"
	"go-mmcli-svr/backend/mmcli"
	_ "go-mmcli-svr/docs"
	"go-mmcli-svr/middleware"
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

		registerModemRoutes(r, backend)
		registerLocationRoutes(r, backend)
		registerSmsRoutes(r, backend)

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

func registerModemRoutes(r chi.Router, backend be.Backend) {
	a := api.Provide(backend)
	r.Get("/modem/", a.ModemList)
	r.Get("/modem/{id:[a-zA-Z0-9%/]+}", a.ModemDetail)
}

func registerLocationRoutes(r chi.Router, backend be.Backend) {
	a := api.Provide(backend)
	r.Get("/location", a.LocationGet)
}

func registerSmsRoutes(r chi.Router, backend be.Backend) {
	a := api.Provide(backend)
	r.Get("/sms/", a.SmsGet)
	r.Get("/sms/{id:[a-zA-Z0-9%/]+}", a.SmsGet)

	r.Post("/sms", a.SmsSend)

	r.Delete("/sms/{id:[a-zA-Z0-9%/]+}", a.SmsDelete)
}
