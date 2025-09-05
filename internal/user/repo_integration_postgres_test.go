// internal/user/repo_integration_postgres_test.go
package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepo_Create_UniqueEmail_Postgres(t *testing.T) {
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

	host, _ := pgC.Host(ctx)
	port, _ := pgC.MappedPort(ctx, "5432/tcp")

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=postgres sslmode=disable TimeZone=Asia/Bangkok", host, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}))

	repo := NewGormRepository(db)

	u1 := &User{Email: "i@a.com", Name: "Ice", PasswordHash: "hash"}
	require.NoError(t, repo.Create(u1))

	u2 := &User{Email: "i@a.com", Name: "Dup", PasswordHash: "hash2"}
	err = repo.Create(u2)
	require.ErrorIs(t, err, ErrDuplicateEmail)
}
