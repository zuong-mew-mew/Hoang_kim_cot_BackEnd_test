# Tài Liệu API Hệ Thống NowEvent

- **Base URL:** `http://localhost:8080`
- **Content-Type:** `application/json`

---

## 1. Quản Lý Sinh Viên (Users)

### 1.1. Thêm sinh viên
- **Method:** `POST`
- **URL:** `/api/users`
- **Request Body:**
```json
{
  "student_id": "21084201",
  "student_name": "Nguyễn Đăng Dương",
  "gender": "male",
  "email": "duong@iuh.edu.vn",
  "zalo_uid": "zalo_123456",
  "major": "Công nghệ thông tin",
  "batch": 21
}
```

### 1.2. Sửa thông tin sinh viên
- **Method:** `PUT`
- **URL:** `/api/users/{student_id}`
- **Ví dụ:** `http://localhost:8080/api/users/21084201`
- **Request Body:**
```json
{
  "student_name": "Nguyễn Đăng Dương",
  "email": "duong.cntt@iuh.edu.vn",
  "zalo_uid": "zalo_999999",
  "major": "Công nghệ thông tin",
  "batch": 21
}
```

### 1.3. Xóa sinh viên
- **Method:** `DELETE`
- **URL:** `/api/users/{student_id}`
- **Ví dụ:** `http://localhost:8080/api/users/21084201`

---

## 2. Quản Lý Sự Kiện (Events)

### 2.1. Thêm sự kiện
- **Method:** `POST`
- **URL:** `/api/events`
- **Request Body:**
```json
{
  "title": "Chào Tân Sinh Viên K21",
  "content": "Chương trình giao lưu văn nghệ",
  "major": "Công nghệ thông tin",
  "batch": 21,
  "event_time_start": "2026-11-20T08:00:00Z",
  "max_slot": 300,
  "current_slot": 0,
  "status": "Open"
}
```

### 2.2. Xóa sự kiện
- **Method:** `DELETE`
- **URL:** `/api/events/{event_id}`
- **Ví dụ:** `http://localhost:8080/api/events/6738c111222333`

---

## 3. Xem Danh Sách Đăng Ký (Registrations)

### 3.1. Xem danh sách người đã đăng ký theo sự kiện
- **Method:** `GET`
- **URL:** `/api/registrations/event/{event_id}`
- **Ví dụ:** `http://localhost:8080/api/registrations/event/6738c111222333`
