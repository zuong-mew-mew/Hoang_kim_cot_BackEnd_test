package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"hoang_kim_cot/pkg/database"
	"hoang_kim_cot/pkg/model"
)

type Handler struct {
	repo *database.Repository
}

func NewHandler(repo *database.Repository) *Handler {
	return &Handler{repo: repo}
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
