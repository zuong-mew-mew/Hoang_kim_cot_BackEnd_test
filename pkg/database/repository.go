package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"hoang_kim_cot/pkg/model"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(client *mongo.Client, dbName string) *Repository {
	return &Repository{
		db: client.Database(dbName),
	}
}

// ================= USER  =================
// Thêm User
func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	return err
}

// Tìm User theo MSSV
func (r *Repository) GetUserByStudentID(ctx context.Context, studentID string) (*model.User, error) {
	var user model.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"student_id": studentID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Sửa User (Cập nhật tên, email, ngành theo MSSV)
func (r *Repository) UpdateUser(ctx context.Context, studentID string, updateData bson.M) error {
	filter := bson.M{"student_id": studentID}
	update := bson.M{"$set": updateData}
	_, err := r.db.Collection("users").UpdateOne(ctx, filter, update)
	return err
}

// Xóa User theo MSSV
func (r *Repository) DeleteUser(ctx context.Context, studentID string) error {
	filter := bson.M{"student_id": studentID}
	_, err := r.db.Collection("users").DeleteOne(ctx, filter)
	return err
}

// ================= EVENT  =================
// Thêm / Cập nhật Event (Crawler)
func (r *Repository) UpsertEvent(ctx context.Context, event *model.Event) error {
	filter := bson.M{"title": event.Title, "batch": event.Batch}
	update := bson.M{
		"$set": bson.M{
			"content":          event.Content,
			"major":            event.Major,
			"event_time_start": event.EventTimeStart,
			"current_slot":     event.CurrentSlot,
			"max_slot":         event.MaxSlot,
			"status":           event.Status,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.db.Collection("events").UpdateOne(ctx, filter, update, opts)
	return err
}

// Xóa Event theo ID
func (r *Repository) DeleteEvent(ctx context.Context, eventID primitive.ObjectID) error {
	filter := bson.M{"_id": eventID}
	_, err := r.db.Collection("events").DeleteOne(ctx, filter)
	return err
}

// Lấy danh sách Event theo Khoa & Khóa
func (r *Repository) GetEventsByMajorAndBatch(ctx context.Context, major string, batch int) ([]model.Event, error) {
	filter := bson.M{
		"major":  major,
		"batch":  batch,
		"status": "Open",
	}

	cursor, err := r.db.Collection("events").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []model.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// ============== Này bên Event ủa lộn registration ==========
// Tạo vé đăng ký
func (r *Repository) CreateRegistration(ctx context.Context, reg *model.Registration) error {
	reg.RegisteredAt = time.Now()
	_, err := r.db.Collection("registrations").InsertOne(ctx, reg)
	return err
}

// Lấy tất cả người đã đăng ký  1 sự kiện nào đó
func (r *Repository) GetRegistrationsByEventID(ctx context.Context, eventID primitive.ObjectID) ([]model.Registration, error) {
	filter := bson.M{"event_id": eventID}

	cursor, err := r.db.Collection("registrations").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Registration
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}
