package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID    string             `bson:"student_id" json:"student_id"`
	StudentName  string             `bson:"student_name" json:"student_name"`
	PasswordHash string             `bson:"password_hash" json:"_"`
	Gender       string             `bson:"gender" json:"gender"`
	Email        string             `bson:"email" json:"email"`
	ZaloUID      string             `bson:"zalo_uid,omitempty" json:"zalo_uid"`
	Major        string             `bson:"major" json:"major"`
	Batch        int                `bson:"batch" json:"batch"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
