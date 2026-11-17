package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationLog: 1 lần gửi thông báo — khớp type NotificationLog của FE.
type NotificationLog struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EventID           *primitive.ObjectID `bson:"event_id,omitempty" json:"eventId,omitempty"`
	Title             string             `bson:"title" json:"title"`
	Body              string             `bson:"body" json:"body"`
	Condition         string             `bson:"condition" json:"condition"`
	TargetDescription string             `bson:"target_description" json:"targetDescription"`
	Status            string             `bson:"status" json:"status"` // "sent" | "failed"
	SentAt            time.Time          `bson:"sent_at" json:"sentAt"`
	SuccessCount      *int               `bson:"success_count,omitempty" json:"successCount,omitempty"`
	FailureCount      *int               `bson:"failure_count,omitempty" json:"failureCount,omitempty"`
	Error             string             `bson:"error,omitempty" json:"error,omitempty"`
	DurationMs 		  int64  			 `bson:"duration_ms,omitempty" json:"durationMs,omitempty"`
}
