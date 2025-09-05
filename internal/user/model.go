package user

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	Name         string    `gorm:"size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
