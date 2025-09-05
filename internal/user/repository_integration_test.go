package user

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func newTestDB(t *testing.T) *gorm.DB {
// 	t.Helper()
// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	require.NoError(t, err)
// 	require.NoError(t, db.AutoMigrate(&User{}))
// 	return db
// }

// func TestRepo_Create_And_UniqueEmail(t *testing.T) {
// 	db := newTestDB(t)
// 	repo := NewGormRepository(db)

// 	u1 := &User{Email: "i@a.com", Name: "Ice", PasswordHash: "hash"}
// 	require.NoError(t, repo.Create(u1))
// 	require.NotZero(t, u1.ID)

// 	u2 := &User{Email: "i@a.com", Name: "Dup", PasswordHash: "hash2"}
// 	err := repo.Create(u2)
// 	require.Error(t, err)
// 	require.ErrorIs(t, err, ErrDuplicateEmail)
// }
