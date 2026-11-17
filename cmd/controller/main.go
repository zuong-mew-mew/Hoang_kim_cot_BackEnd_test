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
)

func withCORS(next http.Handler) http.Handler {
	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000" // mặc định cho local
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	for _, path := range []string{".env", "../../.env"} {
		_ = godotenv.Load(path)
	}

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

	// FCM: thiếu file key thì BE vẫn chạy, chỉ không gửi được (endpoint báo 503).
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
		fmt.Fprintf(w, `{"status": "success", "message": " Server Hoàng Kim Cốt đang chạy bằng net/http!"}`)
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

	fmt.Println(" Server net/http đang chạy tại: http://localhost:8080")
	if err := http.ListenAndServe(":8080", withCORS(mux)); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}
