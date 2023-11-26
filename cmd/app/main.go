package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/commedesvlados/url-shortener/internal/config"
	"github.com/commedesvlados/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/commedesvlados/url-shortener/internal/lib/logger/sl"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/redirect"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/url/remove"
	"github.com/commedesvlados/url-shortener/internal/server/handlers/url/save"
	mwLogger "github.com/commedesvlados/url-shortener/internal/server/middleware/logger"
	"github.com/commedesvlados/url-shortener/internal/storage/sqlite"
)

func main() {
	config.MustLoadVariables()

	log := setupLogger(config.C.Env)

	log.Info("[App] Application started", slog.String("env", config.C.Env), slog.String("version", "1.0"))
	log.Debug("debug are enabled")

	storage, err := sqlite.New(config.E.Database.Path)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	//router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			config.E.HTTPServer.User: config.E.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", remove.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	// Run Server
	log.Info("starting server", slog.String("address", config.E.HTTPServer.Address))

	srv := &http.Server{
		Addr:         config.E.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  config.E.HTTPServer.Timeout,
		WriteTimeout: config.E.HTTPServer.Timeout,
		IdleTimeout:  config.E.HTTPServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvironmentProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case config.EnvironmentDevelopment:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvironmentLocal:
		log = setupPrettySlog()
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
