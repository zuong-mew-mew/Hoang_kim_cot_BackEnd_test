package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Event đại diện cho    {1 sự kiện}   trong hệ thống
type Event struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	Major          string             `bson:"major" json:"major"`
	Batch          int                `bson:"batch" json:"batch"`
	EventTimeStart time.Time          `bson:"event_time_start" json:"event_time_start"`
	CurrentSlot    int                `bson:"current_slot" json:"current_slot"`
	MaxSlot        int                `bson:"max_slot" json:"max_slot"`
	Status         string             `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
}
