package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"hoang_kim_cot/pkg/api"
	"hoang_kim_cot/pkg/database"
	"hoang_kim_cot/pkg/notify"

	"os"
	"github.com/joho/godotenv"
	"strings"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed := strings.Split(os.Getenv("CORS_ORIGIN"), ",")
		if len(allowed) == 0 || allowed[0] == "" {
			allowed = []string{"http://localhost:3000"}
		}
		origin := r.Header.Get("Origin")
		for _, a := range allowed {
			if strings.TrimSpace(a) == origin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load(".env", "../../.env")

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI chưa được cấu hình")
	}

	client, err := database.ConnectMongoDB(uri)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB: %v", err)
	}
	defer client.Disconnect(context.Background())

	repo := database.NewRepository(client, "iuh_event_db")
	handler := api.NewHandler(repo)

	ctx := context.Background()

	if notifyClient, err := notify.New(ctx, os.Getenv("FIREBASE_CREDENTIALS")); err != nil {
		log.Printf("FCM init lỗi: %v", err)
	} else {
		handler.SetNotifyClient(notifyClient)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok"}`)
	})

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "success", "message": "Server Hoàng Kim Cốt đang chạy bằng net/http!"}`)
	})

	mux.HandleFunc("POST /api/auth/login", handler.Login)

	mux.HandleFunc("PUT /api/users/{student_id}", handler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{student_id}", handler.DeleteUser)
	mux.HandleFunc("DELETE /api/events/{event_id}", handler.DeleteEvent)

	mux.HandleFunc("GET /api/users", handler.GetUsers)
	mux.HandleFunc("POST /api/users", handler.CreateUser)
	mux.HandleFunc("GET /api/users/{student_id}", handler.GetUser)

	mux.HandleFunc("GET /api/topics", handler.GetTopics)
	mux.HandleFunc("POST /api/events", handler.CreateEvent)
	mux.HandleFunc("GET /api/events", handler.GetEvents)

	mux.HandleFunc("POST /api/registrations", handler.RegisterEvent)
	mux.HandleFunc("GET /api/registrations/event/{event_id}", handler.GetRegistrationsByEvent)

	mux.HandleFunc("POST /api/notifications/send", handler.SendNotification)
	mux.HandleFunc("GET /api/notifications", handler.GetNotifications)
	mux.HandleFunc("POST /api/users/{student_id}/fcm-token", handler.UpdateFCMToken)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	fmt.Printf("Server net/http đang chạy tại: %s\n", addr)

	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}
