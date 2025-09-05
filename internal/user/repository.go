package user

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrDuplicateEmail = errors.New("email already exists")

type Repository interface {
	Create(u *User) error
	FindByEmail(email string) (*User, error)
}

type GormRepository struct{ db *gorm.DB }

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(u *User) error {
	if err := r.db.Create(u).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *GormRepository) FindByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// func isUniqueErr(err error) bool {
// 	// สำหรับตัวอย่าง: พอเจอ "UNIQUE" ให้ถือว่า duplicate
// 	return err != nil && (contains(err.Error(), "UNIQUE") || contains(err.Error(), "unique"))
// }

// func contains(s, sub string) bool {
// 	return len(s) >= len(sub) && (len(sub) == 0 || (func() bool {
// 		for i := 0; i+len(sub) <= len(s); i++ {
// 			if s[i:i+len(sub)] == sub {
// 				return true
// 			}
// 		}
// 		return false
// 	})())
// }
