package user

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

type fakeSvc struct {
	create func(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
}

func (f *fakeSvc) Create(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	return f.create(ctx, req)
}

func TestHandler_CreateUser_Success(t *testing.T) {
	app := fiber.New()
	svc := &fakeSvc{
		create: func(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
			return &CreateUserResponse{
				ID:        1,
				Email:     req.Email,
				Name:      req.Name,
				CreatedAt: time.Now(),
			}, nil
		},
	}
	RegisterRoutes(app, svc)

	body := []byte(`{"email":"ok@ex.com","name":"Ok","password":"verystrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)
}

func TestHandler_CreateUser_BadRequest(t *testing.T) {
	app := fiber.New()
	svc := &fakeSvc{
		create: func(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
			return nil, ErrBadEmail
		},
	}
	RegisterRoutes(app, svc)

	body := []byte(`{"email":"bad","name":"x","password":"12345678"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}
