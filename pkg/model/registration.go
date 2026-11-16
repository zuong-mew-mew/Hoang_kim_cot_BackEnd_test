package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Registration là chiếc vé ghi nhận 1 lượt đăng ký sự kiện của sinh viên cụ thể
type Registration struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	StudentID    string             `bson:"student_id" json:"student_id"`
	EventID      primitive.ObjectID `bson:"event_id" json:"event_id"`
	RegisteredAt time.Time          `bson:"registered_at" json:"registered_at"`
	Status       string             `bson:"status" json:"status"`
}
