package handlers

import (
	"github.com/badoux/checkmail"
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

var (
	db = database.DB
)

func CreateUser(c *fiber.Ctx) error {
	json := new(structur.CreateUserRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	password := helper.HashAndSalt([]byte(json.Password))
	err := checkmail.ValidateFormat(json.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Email Address",
		})
	}

	new := User{
		Username: json.Username,
		Password: password,
		Email:    json.Email,
		ID:       guuid.New(),
	}

	found := User{}
	query := User{Username: json.Username}
	err = db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return c.Status(400).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	query = User{Email: json.Email}
	err = db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return c.Status(400).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	err = db.Create(&new).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Invalid query database",
			"details": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Succes create user",
	})
}

func GetUsers(c *fiber.Ctx) error {
	User := []User{}
	db.Model(&model.User{}).Order("ID DESC").Find(User)
	return c.Status(200).JSON(fiber.Map{
		"message": "Succes get data",
		"data":    User,
	})
}

func GetUserByUsername(c *fiber.Ctx) error {
	email := c.Params("email")

	user := User{}
	query := User{Email: email}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data":    user,
		"message": "Success get users",
	})

}

func UpdateUser(c *fiber.Ctx) error {
	json := new(structur.CreateUserRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}
	found := User{}
	query := User{Username: json.Username}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "Username not found",
		})
	}

	new := User{
		Username: json.Username,
		Email:    json.Email,
	}

	if json.Email == "" {
		new.Email = found.Email
	}

	if json.Username == "" {
		new.Username = found.Username
	}

	db.Save(&new)
	return c.Status(200).JSON(fiber.Map{
		"message": "Sucess update user",
	})

}

func ChangePassword(c *fiber.Ctx) error {
	json := new(structur.ChangePasswordRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	user := User{}

	if !helper.ComparePasswords(user.Password, []byte(json.Password)) {
		return c.Status(4001).JSON(fiber.Map{
			"message": "Invalid Password",
		})
	}

	user.Password = helper.HashAndSalt([]byte(json.NewPassword))
	db.Save(&user)
	return c.Status(200).JSON(fiber.Map{
		"message": "Sucess update password",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB
	json := new(structur.DeleteUserRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	user := User{}
	// query := User{Username: json.Username, Role: "ROLE-1"}
	query := User{Username: json.Username}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "No user found with given username.",
		})
	}
	if !helper.ComparePasswords(user.Password, []byte(json.Password)) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid Password",
		})
	}

	db.Delete(&user)
	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	c.ClearCookie("sessionid")
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "sucess",
	})
}
