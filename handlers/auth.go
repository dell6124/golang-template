package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

type User model.User
type Session model.Session
type Product model.Product

var (
	SecretKey = []byte("secretminirciuyernmcdkowoir8u{?|>LJNXEkak,}")
)

func GetUser(sessionid guuid.UUID) (User, error) {
	db := database.DB
	query := Session{Sessionid: sessionid}
	found := Session{}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, err
	}
	user := User{}
	usrQuery := User{ID: found.UserRefer}
	err = db.First(&user, &usrQuery).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, err
	}
	return user, nil
}

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	db := database.DB
	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}
	if !helper.ComparePasswords(found.Password, []byte(json.Password)) {
		return c.JSON(fiber.Map{
			"code":    401,
			"message": "Invalid Password",
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  found.Username,
		"admin": true,
		"exp":   time.Now().Add(time.Minute * 1).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(SecretKey)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Invalid generate token",
			"details": err,
		})
	}

	session := Session{
		Sessionid: guuid.New(),
		Expires:   time.Now().Add(time.Minute * 5),
		UserRefer: found.ID,
	}

	err = db.Create(&session).Error // Save session to database
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Invalid query database",
			"details": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"token":   t,
	})
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}
	session := Session{}
	query := Session{Sessionid: json.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "Session not found",
		})
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
}

func GetUserInfo(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "Hallo user",
	})
}
