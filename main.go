package main

import (
	"automated-golang/internal/user"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=Asia/Bangkok",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"postgres",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal(err)
	}

	repo := user.NewGormRepository(db)
	svc := user.NewService(repo)

	app := fiber.New()
	user.RegisterRoutes(app, svc)

	log.Println("listening on :8080")
	log.Fatal(app.Listen(":8080"))
}
