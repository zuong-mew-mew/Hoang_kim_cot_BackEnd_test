package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"hoang_kim_cot/pkg/api"
	"hoang_kim_cot/pkg/database"
)

func main() {
	uri := ""

	client, err := database.ConnectMongoDB(uri)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB: %v", err)
	}
	defer client.Disconnect(context.Background())

	repo := database.NewRepository(client, "iuh_event_db")
	handler := api.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "success", "message": " Server Hoàng Kim Cốt đang chạy bằng net/http!"}`)
	})

	mux.HandleFunc("PUT /api/users/{student_id}", handler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{student_id}", handler.DeleteUser)
	mux.HandleFunc("DELETE /api/events/{event_id}", handler.DeleteEvent)

	mux.HandleFunc("POST /api/users", handler.CreateUser)
	mux.HandleFunc("GET /api/users/{student_id}", handler.GetUser)

	mux.HandleFunc("POST /api/events", handler.CreateEvent)
	mux.HandleFunc("GET /api/events", handler.GetEvents)

	mux.HandleFunc("POST /api/registrations", handler.RegisterEvent)
	mux.HandleFunc("GET /api/registrations/event/{event_id}", handler.GetRegistrationsByEvent)

	fmt.Println(" Server net/http đang chạy tại: http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}
