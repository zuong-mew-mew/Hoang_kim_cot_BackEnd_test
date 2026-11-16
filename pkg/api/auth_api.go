package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// 1. API: Thêm sinh viên mới (POST /api/users)
func (h *Handler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu gửi lên không hợp lệ"})
		return
	}

	err := h.repo.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo sinh viên"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tạo sinh viên thành công!", "data": user})
}

// 2. API: Lấy thông tin sinh viên theo MSSV (GET /api/users/:student_id)
func (h *Handler) GetUser(c *gin.Context) {
	studentID := c.Param("student_id")

	user, err := h.repo.GetUserByStudentID(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sinh viên"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// 3. API: Tạo sự kiện mới (POST /api/events)
func (h *Handler) CreateEvent(c *gin.Context) {
	var event model.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu sự kiện không hợp lệ"})
		return
	}

	err := h.repo.UpsertEvent(c.Request.Context(), &event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo sự kiện"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lưu sự kiện thành công!", "data": event})
}

// 4. API: Lấy sự kiện theo Khoa & Khóa (GET /api/events?major=...&batch=...)
func (h *Handler) GetEvents(c *gin.Context) {
	major := c.Query("major")
	batchStr := c.Query("batch")
	batch, _ := strconv.Atoi(batchStr)

	events, err := h.repo.GetEventsByMajorAndBatch(c.Request.Context(), major, batch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy sự kiện"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// 5. API: Sinh viên bấm Đăng ký vé (POST /api/registrations)
func (h *Handler) RegisterEvent(c *gin.Context) {
	var reg model.Registration
	if err := c.ShouldBindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đăng ký không hợp lệ"})
		return
	}

	err := h.repo.CreateRegistration(c.Request.Context(), &reg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Đăng ký thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đăng ký thành công!", "data": reg})
}

// 6. API: Xem danh sách những ai đã đăng ký sự kiện này (GET /api/registrations/event/:event_id)
func (h *Handler) GetRegistrationsByEvent(c *gin.Context) {
	eventIDStr := c.Param("event_id")
	eventID, err := primitive.ObjectIDFromHex(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sự kiện không đúng định dạng"})
		return
	}

	list, err := h.repo.GetRegistrationsByEventID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách"})
		return
	}

	c.JSON(http.StatusOK, list)
}
