package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"hoang_kim_cot/pkg/database"
	"hoang_kim_cot/pkg/model"
	"hoang_kim_cot/pkg/topics"
	"hoang_kim_cot/pkg/notify"

	"crypto/rand"
	"encoding/hex"
	"os"

	"fmt"
	"sort"
	"strings"
	"time"
)

type Handler struct {
	repo *database.Repository
	notify *notify.Client
}

func NewHandler(repo *database.Repository) *Handler {
	return &Handler{repo: repo}
}

// POST /api/auth/login — chỉ admin đăng nhập dashboard
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dữ liệu JSON không hợp lệ"})
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminEmail == "" || adminPassword == "" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Chưa cấu hình tài khoản admin"})
		return
	}

	if req.Email != adminEmail || req.Password != adminPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Sai email hoặc mật khẩu"})
		return
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không tạo được phiên đăng nhập"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": hex.EncodeToString(token),
		"user": map[string]string{
			"name":  "Quản trị viên",
			"email": adminEmail,
			"role":  "admin",
		},
	})
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.repo.GetAllUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lỗi lấy danh sách sinh viên"})
		return
	}
	if users == nil {
		users = []model.User{} // trả [] thay vì null
	}
	json.NewEncoder(w).Encode(users)
}

// 1. POST /api/users (Thêm sinh viên)
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dữ liệu JSON không hợp lệ"})
		return
	}

	err := h.repo.CreateUser(r.Context(), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không thể tạo sinh viên"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tạo sinh viên thành công!",
		"data":    user,
	})
}

// 2. GET /api/users/{student_id} (Tìm sinh viên theo MSSV)
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentID := r.PathValue("student_id")

	user, err := h.repo.GetUserByStudentID(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không tìm thấy sinh viên"})
		return
	}

	json.NewEncoder(w).Encode(user)
}

// 3. POST /api/events (Tạo sự kiện)
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dữ liệu sự kiện không hợp lệ"})
		return
	}

	err := h.repo.UpsertEvent(r.Context(), &event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không thể lưu sự kiện"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Lưu sự kiện thành công!",
		"data":    event,
	})
}

// 4. GET /api/events?major=...&batch=... (Lấy sự kiện)
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	major := r.URL.Query().Get("major")
	batchStr := r.URL.Query().Get("batch")
	batch, _ := strconv.Atoi(batchStr)

	events, err := h.repo.GetEventsByMajorAndBatch(r.Context(), major, batch)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lỗi lấy sự kiện"})
		return
	}

	json.NewEncoder(w).Encode(events)
}

// 5. POST /api/registrations (Đăng ký vé)
func (h *Handler) RegisterEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reg model.Registration
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dữ liệu đăng ký không hợp lệ"})
		return
	}

	err := h.repo.CreateRegistration(r.Context(), &reg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Đăng ký thất bại"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Đăng ký thành công!",
		"data":    reg,
	})
}

// 6. GET /api/registrations/event/{event_id} (Xem danh sách vé)
func (h *Handler) GetRegistrationsByEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	eventIDStr := r.PathValue("event_id")
	eventID, err := primitive.ObjectIDFromHex(eventIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID sự kiện không đúng định dạng"})
		return
	}

	list, err := h.repo.GetRegistrationsByEventID(r.Context(), eventID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lỗi lấy danh sách"})
		return
	}

	json.NewEncoder(w).Encode(list)
}

// PUT /api/users/{student_id} (Sửa thông tin sinh viên)
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentID := r.PathValue("student_id")

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dữ liệu JSON không hợp lệ"})
		return
	}

	err := h.repo.UpdateUser(r.Context(), studentID, updateData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không thể cập nhật sinh viên"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Cập nhật sinh viên thành công!"})
}

// DELETE /api/users/{student_id} (Xóa sinh viên)
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentID := r.PathValue("student_id")

	err := h.repo.DeleteUser(r.Context(), studentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không thể xóa sinh viên"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Xóa sinh viên thành công!"})
}

// DELETE /api/events/{event_id} (Xóa sự kiện)
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	eventIDStr := r.PathValue("event_id")

	eventID, err := primitive.ObjectIDFromHex(eventIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID sự kiện không hợp lệ"})
		return
	}

	err = h.repo.DeleteEvent(r.Context(), eventID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không thể xóa sự kiện"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Xóa sự kiện thành công!"})
}


// Topic: 1 topic FCM trong dashboard — khớp đúng type Topic của FE.
type Topic struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Kind            string `json:"kind"` // "faculty" | "cohort"
	Label           string `json:"label"`
	Description     string `json:"description,omitempty"`
	SubscriberCount int    `json:"subscriberCount"`
	Enabled         bool   `json:"enabled"`
}

// GET /api/topics — mỗi major = 1 topic khoa, mỗi batch = 1 topic khóa.
func (h *Handler) GetTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	byMajor, byBatch, err := h.repo.GetTopicStats(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lỗi lấy danh sách topic"})
		return
	}

	topicsList := []Topic{}

	// Topic khoa: tên topic tra từ package topics, major chưa có mapping
	// thì slug hoá từ tên để không mất dữ liệu.
	for _, major := range sortedKeys(byMajor) {
		name, ok := topics.ResolveFacultyTopic(major)
		if !ok {
			name = "faculty_" + strings.ReplaceAll(topics.Normalize(major), " ", "_")
		}
		topicsList = append(topicsList, Topic{
			ID: name, Name: name, Kind: "faculty",
			Label:       major,
			Description: "Sinh viên ngành " + major,
			SubscriberCount: byMajor[major],
			Enabled:     true,
		})
	}

	// Topic khóa: batch 21 -> cohort_k21.
	for _, batch := range sortedIntKeys(byBatch) {
		name := topics.ResolveCohortTopic(batch)
		topicsList = append(topicsList, Topic{
			ID: name, Name: name, Kind: "cohort",
			Label:       fmt.Sprintf("Khóa K%d", batch),
			Description: fmt.Sprintf("Sinh viên khóa K%d", batch),
			SubscriberCount: byBatch[batch],
			Enabled:     true,
		})
	}

	if topicsList == nil {
		topicsList = []Topic{}
	}
	json.NewEncoder(w).Encode(topicsList)
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedIntKeys(m map[int]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}


// SetNotifyClient gắn FCM client — gọi sau NewHandler; không gọi thì gửi báo 503.
func (h *Handler) SetNotifyClient(c *notify.Client) {
	h.notify = c
}

type sendNotificationRequest struct {
	EventID   string   `json:"event_id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Faculty   string   `json:"faculty"`
	Cohorts   []string `json:"cohorts"`
	Condition string   `json:"condition"`
}

// POST /api/notifications/send — gửi FCM rồi ghi log.
// Có event_id thì BE tự tính condition từ major/batch của event;
// không thì dùng condition FE gửi lên.
func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req sendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dữ liệu không hợp lệ"})
		return
	}
	if req.Title == "" || req.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Thiếu tiêu đề hoặc nội dung"})
		return
	}

	logEntry := model.NotificationLog{
		Title: req.Title, Body: req.Body, Status: "failed",
	}

	data := map[string]string{}
	tokens := []string{}

	if req.EventID != "" {
		eventID, err := primitive.ObjectIDFromHex(req.EventID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Event ID không hợp lệ"})
			return
		}
		event, err := h.repo.GetEventByID(r.Context(), eventID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Không tìm thấy sự kiện"})
			return
		}
		logEntry.EventID = &eventID
		facultyTopic, _ := topics.ResolveFacultyTopic(event.Major)
		cohortTopic := topics.ResolveCohortTopic(event.Batch)
		logEntry.Condition = topics.BuildCondition(facultyTopic, cohortTopic)
		logEntry.TargetDescription = fmt.Sprintf("%s — Khóa K%d", event.Major, event.Batch)
		data["event_id"] = req.EventID
		tokens, _ = h.repo.GetFCMTokensByMajorAndBatch(r.Context(), event.Major, event.Batch)
	} else {
		if req.Condition == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Thiếu condition gửi FCM"})
			return
		}
		logEntry.Condition = req.Condition
		logEntry.TargetDescription = req.Condition
	}

	if logEntry.Condition == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không tính được condition gửi FCM"})
		return
	}
	if h.notify == nil || !h.notify.Ready() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "FCM chưa được cấu hình trên BE"})
		return
	}

	sendStart := time.Now()
	if len(tokens) > 0 {
		success, failure, err := h.notify.SendMulticast(r.Context(), tokens, req.Title, req.Body, data)
		logEntry.DurationMs = time.Since(sendStart).Milliseconds()
		if err != nil {
			logEntry.Error = err.Error()
			_ = h.repo.SaveNotificationLog(r.Context(), &logEntry)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Gửi FCM thất bại: " + err.Error()})
			return
		}
		logEntry.SuccessCount = &success
		logEntry.FailureCount = &failure
	} else {
		if _, err := h.notify.SendToCondition(r.Context(), logEntry.Condition, req.Title, req.Body, data); err != nil {
			logEntry.DurationMs = time.Since(sendStart).Milliseconds()
			logEntry.Error = err.Error()
			_ = h.repo.SaveNotificationLog(r.Context(), &logEntry)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Gửi FCM thất bại: " + err.Error()})
			return
		}
		logEntry.DurationMs = time.Since(sendStart).Milliseconds()
	}

	logEntry.Status = "sent"
	if err := h.repo.SaveNotificationLog(r.Context(), &logEntry); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Gửi được nhưng không lưu được log"})
		return
	}
	json.NewEncoder(w).Encode(logEntry)
}

// GET /api/notifications — 100 log mới nhất.
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logs, err := h.repo.GetNotifications(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lỗi lấy lịch sử thông báo"})
		return
	}
	if logs == nil {
		logs = []model.NotificationLog{}
	}
	json.NewEncoder(w).Encode(logs)
}

// POST /api/users/{student_id}/fcm-token — app báo token mới lên.
func (h *Handler) UpdateFCMToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	studentID := r.PathValue("student_id")
	var req struct {
		FCMToken string `json:"fcm_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.FCMToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Thiếu fcm_token"})
		return
	}
	if err := h.repo.UpdateUserFCMToken(r.Context(), studentID, req.FCMToken); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Không lưu được token"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Đã lưu token"})
}
