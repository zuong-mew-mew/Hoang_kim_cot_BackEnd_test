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

func (r *Repository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	cursor, err := r.db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := []model.User{}
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
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
	filter := bson.M{}
	// Chỉ lọc theo major/batch khi có giá trị thật — không thì lấy tất cả.
	if major != "" {
		filter["major"] = major
	}
	if batch != 0 {
		filter["batch"] = batch
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

// GetTopicStats đếm số sinh viên theo từng major và từng batch.
func (r *Repository) GetTopicStats(ctx context.Context) (map[string]int, map[int]int, error) {
	users, err := r.GetAllUsers(ctx)
	if err != nil {
		return nil, nil, err
	}
	byMajor := map[string]int{}
	byBatch := map[int]int{}
	for _, u := range users {
		if u.Major != "" {
			byMajor[u.Major]++
		}
		if u.Batch != 0 {
			byBatch[u.Batch]++
		}
	}
	return byMajor, byBatch, nil
}


// Lưu 1 log thông báo vào collection "notifications".
func (r *Repository) SaveNotificationLog(ctx context.Context, log *model.NotificationLog) error {
	log.SentAt = time.Now()
	res, err := r.db.Collection("notifications").InsertOne(ctx, log)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}
	return nil
}

// GetNotifications trả log mới nhất trước, tối đa 100 dòng.
func (r *Repository) GetNotifications(ctx context.Context) ([]model.NotificationLog, error) {
	opts := options.Find().SetSort(bson.D{{Key: "sent_at", Value: -1}}).SetLimit(100)
	cursor, err := r.db.Collection("notifications").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	logs := []model.NotificationLog{}
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GetEventByID lấy 1 event theo _id (để gửi thông báo theo event).
func (r *Repository) GetEventByID(ctx context.Context, id primitive.ObjectID) (*model.Event, error) {
	var event model.Event
	if err := r.db.Collection("events").FindOne(ctx, bson.M{"_id": id}).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

// UpdateUserFCMToken lưu FCM token của 1 sinh viên (app gửi lên mỗi lần
// đăng nhập và mỗi khi token đổi).
func (r *Repository) UpdateUserFCMToken(ctx context.Context, studentID string, token string) error {
	filter := bson.M{"student_id": studentID}
	update := bson.M{"$set": bson.M{"fcm_token": token}}
	_, err := r.db.Collection("users").UpdateOne(ctx, filter, update)
	return err
}

// GetFCMTokensByMajorAndBatch trả token của các user khớp khoa + khóa,
// chỉ lấy những user đã báo token lên.
func (r *Repository) GetFCMTokensByMajorAndBatch(ctx context.Context, major string, batch int) ([]string, error) {
	filter := bson.M{
		"major":     major,
		"batch":     batch,
		"fcm_token": bson.M{"$ne": ""},
	}
	cursor, err := r.db.Collection("users").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	tokens := []string{}
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		if user.FCMToken != "" {
			tokens = append(tokens, user.FCMToken)
		}
	}
	return tokens, cursor.Err()
}
