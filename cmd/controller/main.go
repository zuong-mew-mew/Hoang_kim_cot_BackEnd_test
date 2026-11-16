package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"hoang_kim_cot/pkg/api"
	"hoang_kim_cot/pkg/database"
)

func main() {
	uri := "mongodb+srv://gaulmtvn77_db_user:Gaulmtvn2026@cluster0.sqbjtgf.mongodb.net/?retryWrites=true&w=majority"

	// 1. Kết nối MongoDB
	client, err := database.ConnectMongoDB(uri)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// 2. Khởi tạo Repo và Handler
	repo := database.NewRepository(client, "iuh_event_db")
	handler := api.NewHandler(repo)

	// 3. Khởi tạo Gin Router
	r := gin.Default()

	// 4. Đăng ký các đường link API
	r.POST("/api/users", handler.CreateUser)
	r.GET("/api/users/:student_id", handler.GetUser)

	r.POST("/api/events", handler.CreateEvent)
	r.GET("/api/events", handler.GetEvents)

	r.POST("/api/registrations", handler.RegisterEvent)
	r.GET("/api/registrations/event/:event_id", handler.GetRegistrationsByEvent)

	// 5. Bật server chạy tại cổng 8080
	fmt.Println(" Server đang chạy tại: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Lỗi khởi động server: %v", err)
	}
}
