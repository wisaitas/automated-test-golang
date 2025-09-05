// internal/user/e2e_test.go
package user

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func buildApp(t *testing.T) *fiber.App {
	t.Helper()

	// Spin up a disposable Postgres
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17",
		Env:          map[string]string{"POSTGRES_USER": "postgres", "POSTGRES_PASSWORD": "postgres", "POSTGRES_DB": "postgres"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgC.Terminate(ctx) })

	host, err := pgC.Host(ctx)
	require.NoError(t, err)
	mapped, err := pgC.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=postgres sslmode=disable TimeZone=Asia/Bangkok", host, mapped.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}))

	repo := NewGormRepository(db)
	svc := NewService(repo)

	app := fiber.New()
	RegisterRoutes(app, svc)
	return app
}

func TestE2E_CreateUser_Then_Duplicate(t *testing.T) {
	app := buildApp(t)

	// ครั้งที่ 1: 201
	body1 := []byte(`{"email":"e2e@ok.com","name":"E2E","password":"verystrong"}`)
	req1 := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")
	res1, err := app.Test(req1)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res1.StatusCode)

	// ครั้งที่ 2 (email ซ้ำ): 409
	body2 := []byte(`{"email":"e2e@ok.com","name":"Dup","password":"verystrong"}`)
	req2 := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	res2, err := app.Test(req2)
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, res2.StatusCode)
}
