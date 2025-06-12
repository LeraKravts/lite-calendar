package main

import (
	"fmt"

	"github.com/lerakravts/lite-calendar/internal/handler"
	"github.com/lerakravts/lite-calendar/internal/repo"
	"github.com/lerakravts/lite-calendar/internal/service"

	"log"
	"net/http"
)

func main() {

	//подключаем все слои и создаем зависимости
	repo := repo.NewRepository()
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	// Create route dispatcher
	mux := http.NewServeMux()

	//   Обрабатываем оба метода (GET и POST) по пути /events:
	//валидация метода происходит во внешнем switch, прямо в main.go.
	//
	//А значит, GetEvents и CreateEvent не обязаны проверять метод внутри себя
	//— они вызываются только если метод уже прошёл проверку.
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEvents(w, r)
		case http.MethodPost:
			handler.CreateEvent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/events/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler.DeleteEvent(w, r)
			return

		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// Start server (listens on port 8080 and routes to mux)
	fmt.Println("🚀 Starting server at http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("❌ Server failed to start", err)
	}

}
