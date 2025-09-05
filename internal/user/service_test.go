package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	createFn      func(u *User) error
	findByEmailFn func(email string) (*User, error)
}

func (m *mockRepo) Create(u *User) error                    { return m.createFn(u) }
func (m *mockRepo) FindByEmail(email string) (*User, error) { return m.findByEmailFn(email) }

func TestServiceCreate_Success(t *testing.T) {
	mr := &mockRepo{
		createFn: func(u *User) error {
			u.ID = 123
			u.CreatedAt = time.Now()
			return nil
		},
		findByEmailFn: func(email string) (*User, error) { return nil, nil },
	}
	svc := NewService(mr)

	resp, err := svc.Create(context.Background(), CreateUserRequest{
		Email:    "a@b.com",
		Name:     "Alice",
		Password: "supersecret",
	})
	require.NoError(t, err)
	require.Equal(t, uint(123), resp.ID)
	require.Equal(t, "a@b.com", resp.Email)
	require.Equal(t, "Alice", resp.Name)
}

func TestServiceCreate_InvalidEmail(t *testing.T) {
	mr := &mockRepo{
		createFn:      func(u *User) error { return nil },
		findByEmailFn: func(email string) (*User, error) { return nil, nil },
	}
	svc := NewService(mr)

	_, err := svc.Create(context.Background(), CreateUserRequest{
		Email:    "bad-email",
		Name:     "Alice",
		Password: "supersecret",
	})
	require.ErrorIs(t, err, ErrBadEmail)
}

func TestServiceCreate_Duplicate(t *testing.T) {
	mr := &mockRepo{
		createFn:      func(u *User) error { return ErrDuplicateEmail },
		findByEmailFn: func(email string) (*User, error) { return nil, errors.New("not used") },
	}
	svc := NewService(mr)

	_, err := svc.Create(context.Background(), CreateUserRequest{
		Email:    "dup@b.com",
		Name:     "Bob",
		Password: "password!",
	})
	require.ErrorIs(t, err, ErrDuplicateEmail)
}
