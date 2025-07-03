package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lerakravts/lite-calendar/internal/config"
	"github.com/lerakravts/lite-calendar/internal/db"
	"github.com/lerakravts/lite-calendar/internal/handler"
	"github.com/lerakravts/lite-calendar/internal/logger"
	"github.com/lerakravts/lite-calendar/internal/repo"
	"github.com/lerakravts/lite-calendar/internal/service"
)

func main() {
	cfg := config.Load()
	//Теперь доступно: cfg.AppPort, cfg.PostgresHost, cfg.AppEnv и т.д.

	// 1. Инициализируем логгер
	logger.Init(cfg.AppEnv)
	// 2. Подключаемся к БД
	dbConn, err := db.Connect(cfg)
	if err != nil {
		slog.Error("❌Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbConn.Close()

	slog.Info("📡 Server starting...",
		slog.String("port", cfg.AppPort),
		slog.String("env", cfg.AppEnv),
	)

	// 3. Инициализируем слои и зависимости
	r := repo.NewRepository(dbConn)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	// 🔀 роутинг через chi
	router := chi.NewRouter()
	router.Get("/events", h.GetEvents)
	router.Post("/events", h.CreateEvent)
	router.Delete("/events/{id}", h.DeleteEvent)

	//структура http.Server, которая представляет полноценный HTTP-сервер Go.
	//управление этим сервером напрямую — через srv.ListenAndServe() и srv.Shutdown().
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort, // ← где слушать входящие HTTP-запросы
		Handler: router,            // ← как обрабатывать (наш chi.Router)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		slog.Info("📡 Server starting...",
			slog.String("port", cfg.AppPort),
			slog.String("env", cfg.AppEnv),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("❌ ListenAndServe failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	//ждём сигнала завершения
	<-stop
	slog.Info("⏳ Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("❌ Graceful shutdown failed", slog.String("error", err.Error()))
	} else {
		slog.Info("✅ Server shutdown cleanly")
	}

}
