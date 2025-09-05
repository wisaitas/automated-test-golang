package user

import "errors"

var (
	ErrBadEmail    = errors.New("invalid email")
	ErrBadPassword = errors.New("password must be at least 8 chars")
	ErrBadName     = errors.New("name is required")
)

func ValidateCreate(r CreateUserRequest) error {
	if len(r.Name) == 0 {
		return ErrBadName
	}
	// โค้ดง่ายๆ: ขอมี '@' ก็พอสำหรับตัวอย่าง (จริงควรใช้ lib/email)
	if len(r.Email) < 3 || !containsAt(r.Email) {
		return ErrBadEmail
	}
	if len(r.Password) < 8 {
		return ErrBadPassword
	}
	return nil
}

func containsAt(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}
