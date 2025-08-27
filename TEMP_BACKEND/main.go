package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("supersecretkey")

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Birthday  time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("temp_backend.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	db.AutoMigrate(&User{})
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(u *User) (string, error) {
	claims := jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// @title TEMP_BACKEND API
// @version 1.0
// @description Simple auth example

func main() {
	initDB()
	app := fiber.New()
	app.Use(logger.New())

	app.Post("/register", func(c *fiber.Ctx) error {
		var body struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
			Birthday  string `json:"birthday"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		hashed, err := hashPassword(body.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash"})
		}
		bday, _ := time.Parse("2006-01-02", body.Birthday)
		user := User{Email: body.Email, Password: hashed, FirstName: body.FirstName, LastName: body.LastName, Phone: body.Phone, Birthday: bday}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"id": user.ID, "email": user.Email})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		var user User
		if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		if !checkPassword(user.Password, body.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		tkn, err := generateToken(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate token"})
		}
		return c.JSON(fiber.Map{"token": tkn})
	})

	// protected
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtSecret,
		ContextKey: "user",
	}))

	app.Get("/me", func(c *fiber.Ctx) error {
		userToken := c.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)
		sub := uint(claims["sub"].(float64))
		var user User
		if err := db.First(&user, sub).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		user.Password = ""
		return c.JSON(user)
	})

	app.Listen(":3000")
}
