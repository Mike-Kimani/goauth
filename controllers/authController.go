package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mike-kimani/goauth/database"
	"github.com/mike-kimani/goauth/models"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World👋")
}

const SecretKey = "secret"

var MinPasswordLength = 5

func Register(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)

	if data["name"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	if err != nil {
		return err
	}
	if data["email"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email is required",
		})
	}
	if data["password"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is required",
		})
	}

	if len(data["password"]) < MinPasswordLength {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is too short. Must be at least 5 characters",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}
	database.DB.Create(&user)

	if user.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user not created. check if already exists",
		})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "incorrect password",
		})
		//return c.JSON(fiber.Map{
		//	"message": "incorrect password",
		//})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "server error generating token",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	if user.Id == 0 {
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "successfully logged out",
	})
}
