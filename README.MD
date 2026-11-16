# Tài Liệu API Hệ Thống NowEvent

- **Base URL:** `http://localhost:8080`
- **Content-Type:** `application/json`

---

## 1. Quản Lý Sinh Viên (Users)

### 1.1. Tạo mới sinh viên
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

### 1.2. Tra cứu thông tin sinh viên theo MSSV
- **Method:** `GET`
- **URL:** `/api/users/{student_id}`
- **Ví dụ:** `http://localhost:8080/api/users/21084201`

---

## 2. Quản Lý Sự Kiện (Events)

### 2.1. Tạo mới hoặc cập nhật sự kiện
- **Method:** `POST`
- **URL:** `/api/events`
- **Request Body:**
```json
{
  "title": "Chào Tân Sinh Viên K21",
  "content": "Chương trình giao lưu văn nghệ và định hướng học tập",
  "major": "Công nghệ thông tin",
  "batch": 21,
  "event_time_start": "2026-11-20T08:00:00Z",
  "max_slot": 300,
  "current_slot": 0,
  "status": "Open"
}
```

### 2.2. Lọc danh sách sự kiện theo Khoa và Khóa
- **Method:** `GET`
- **URL:** `/api/events`
- **Query Parameters:**
    - `major`: Tên khoa (ví dụ: `Công nghệ thông tin`)
    - `batch`: Khóa tuyển sinh (ví dụ: `21`)
- **Ví dụ:** `http://localhost:8080/api/events?major=Công nghệ thông tin&batch=21`

---

## 3. Quản Lý Đăng Ký (Registrations)

### 3.1. Đăng ký tham gia sự kiện
- **Method:** `POST`
- **URL:** `/api/registrations`
- **Request Body:**
```json
{
  "user_id": "ID_OBJECT_USER",
  "student_id": "21084201",
  "event_id": "ID_OBJECT_EVENT",
  "status": "SUCCESS"
}
```

### 3.2. Lấy danh sách sinh viên đã đăng ký theo sự kiện
- **Method:** `GET`
- **URL:** `/api/registrations/event/{event_id}`
- **Ví dụ:** `http://localhost:8080/api/registrations/event/6738c111222333`