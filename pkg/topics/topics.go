// Package topics: tính tên topic FCM từ major (tên khoa/ngành) + batch (khóa).
//
// Topic là dữ liệu suy ra từ major + batch — không lưu riêng trong DB.
// Mọi chỗ cần topic (crawler lúc gửi, handler lúc trả user cho app,
// dashboard lúc hiển thị) đều dùng chung package này để không lệch nhau.
package topics

import (
	"fmt"
	"strings"
)

// FacultyMap: tên khoa/ngành cào về (ĐÃ CHUẨN HOÁ) -> slug topic.
// Key phải là tên thật mà web Đoàn trường dùng — bổ sung dần theo log WARN
// khi gặp major chưa có trong bảng.
var FacultyMap = map[string]string{
	"cong nghe thong tin": "cntt",
	// "quan tri kinh doanh": "qtkd",
	// ... 14 khoa còn lại
}

// Normalize: "Khoa Công nghệ   Thông Tin" -> "cong nghe thong tin"
func Normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.TrimPrefix(s, "khoa ")
	var b strings.Builder
	for _, r := range s {
		if rep, ok := vnMap[r]; ok {
			b.WriteRune(rep)
		} else {
			b.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// ResolveFacultyTopic: "Công nghệ thông tin" -> "faculty_cntt", false nếu chưa có mapping
func ResolveFacultyTopic(major string) (string, bool) {
	norm := Normalize(major)
	if slug, ok := FacultyMap[norm]; ok {
		return "faculty_" + slug, true
	}
	// Fallback: khớp chứa — "ky thuat phan mem (cntt)" vẫn bắt được
	for name, slug := range FacultyMap {
		if strings.Contains(norm, name) || strings.Contains(name, norm) {
			return "faculty_" + slug, true
		}
	}
	return "", false
}

// ResolveCohortTopic: 21 -> "cohort_k21" (batch 0 hoặc thiếu -> "" = không giới hạn khóa)
func ResolveCohortTopic(batch int) string {
	if batch <= 0 {
		return ""
	}
	return fmt.Sprintf("cohort_k%02d", batch)
}

// BuildCondition: điều kiện FCM gửi đi
// ("faculty_cntt", "cohort_k21") -> "'faculty_cntt' in topics && 'cohort_k21' in topics"
func BuildCondition(facultyTopic, cohortTopic string) string {
	parts := []string{}
	if facultyTopic != "" {
		parts = append(parts, fmt.Sprintf("'%s' in topics", facultyTopic))
	}
	if cohortTopic != "" {
		parts = append(parts, fmt.Sprintf("'%s' in topics", cohortTopic))
	}
	return strings.Join(parts, " && ")
}

// vnMap: chữ có dấu -> không dấu
var vnMap = map[rune]rune{
	'à': 'a', 'á': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a',
	'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
	'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
	'đ': 'd',
	'è': 'e', 'é': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e',
	'ê': 'e', 'ề': 'e', 'ế': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
	'ì': 'i', 'í': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
	'ò': 'o', 'ó': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o',
	'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
	'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
	'ù': 'u', 'ú': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u',
	'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
	'ỳ': 'y', 'ý': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
}
