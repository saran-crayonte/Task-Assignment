package user

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saran-crayonte/task/database"
	"github.com/saran-crayonte/task/models"
	"golang.org/x/crypto/bcrypt"
)

func Register() fiber.Handler {
	return func(c fiber.Ctx) error {
		dat := new(models.User)
		if err := json.Unmarshal(c.Body(), &dat); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		// Check if username already exists
		var existingUser models.User
		database.DB.First(&existingUser, "username = ?", dat.Username)
		if len(existingUser.Username) != 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "this username already exists"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dat.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "hashing failed"})
		}

		newUser := models.User{
			Username: dat.Username,
			Name:     dat.Name,
			Email:    dat.Email,
			Password: string(hashedPassword),
		}
		database.DB.Create(&newUser)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
		})
	}
}

func Login() fiber.Handler {
	return func(c fiber.Ctx) error {
		returnObject := fiber.Map{
			"status": "",
			"msg":    "Something went wrong.",
		}

		// 1. Check user for the given credentials

		var formData models.User

		if err := json.Unmarshal(c.Body(), &formData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		// log.Println("Form binding error.")
		// return c.Status(fiber.StatusBadRequest).JSON(returnObject)

		var user models.User

		database.DB.First(&user, "username = ?", formData.Username)
		if len(user.Username) == 0 {
			returnObject["msg"] = "User not found."

			return c.Status(fiber.StatusBadRequest).JSON(returnObject)
		}

		// Validate password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(formData.Password))

		if err != nil {
			//log.Println("Invalid password.")

			returnObject["msg"] = "Password doesnt match"
			return c.Status(fiber.StatusUnauthorized).JSON(returnObject)
		}

		// 2. Create token

		token, err := GenerateToken(user)

		if err != nil {
			returnObject["msg"] = "Token creation error."
			return c.Status(fiber.StatusUnauthorized).JSON(returnObject)
		}

		returnObject["token"] = token
		returnObject["user"] = user
		returnObject["status"] = "OK"
		returnObject["msg"] = "User authenticated"
		return c.Status(fiber.StatusAccepted).JSON(returnObject)
	}
}

func UpdatePassword() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the username from the request context
		username, ok := c.Locals("username").(string)
		if !ok {
			return fiber.ErrUnauthorized
		}

		// Parse request body
		var user models.User
		if err := json.Unmarshal(c.Body(), &user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}
		if username != user.Username {
			return fiber.ErrUnauthorized
		}

		var existingUser models.User
		database.DB.First(&existingUser, "username = ?", user.Username)
		if len(existingUser.Username) == 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username doesn't exists"})
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "hashing failed"})
		}

		// Update the user's password in the database
		user.Password = string(hashedPassword)
		database.DB.Model(&existingUser).Updates(user)

		return c.JSON(fiber.Map{
			"message": "Password updated successfully",
		})
	}
}

type CustomClaims struct {
	Email    string
	Username string

	jwt.RegisteredClaims
}

var secret string = "secret"

func GenerateToken(user models.User) (string, error) {

	claims := CustomClaims{
		user.Email,
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Println("Error in token signing.", err)
		return "", err
	}

	return t, nil

}

// Validate Token
func ValidateToken(clientToken string) (*CustomClaims, string) {
	token, err := jwt.ParseWithClaims(clientToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err.Error()
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		return nil, "Invalid token claims"
	}

	return claims, ""
}

func Authenticate() fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("token")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token not present."})
		}

		claims, msg := ValidateToken(token)

		log.Println(claims)

		if msg != "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": msg})
		}

		c.Locals("username", claims.Username)

		return c.Next()
	}
}

func RefreshToken() fiber.Handler {
	return func(c fiber.Ctx) error {
		returnObject := fiber.Map{
			"status": "OK",
			"msg":    "Refresh Token route",
		}

		username := c.Get("username")
		if username == "" {
			log.Println("Username key not found.")

			returnObject["msg"] = "username not found."
			return c.Status(fiber.StatusUnauthorized).JSON(returnObject)
		}

		// emailStr, ok := email.(string)
		// if !ok {
		// 	log.Println("Email is not of type string.")
		// 	returnObject["msg"] = "Email is not valid."
		// 	return c.Status(fiber.StatusBadRequest).JSON(returnObject)
		// }

		var user models.User
		database.DB.First(&user, "username = ?", username)

		if len(user.Username) == 0 {
			returnObject["msg"] = "Username not found."

			return c.Status(fiber.StatusBadRequest).JSON(returnObject)
		}

		token, err := GenerateToken(user)

		if err != nil {
			returnObject["msg"] = "Token creation error."
			return c.Status(fiber.StatusUnauthorized).JSON(returnObject)
		}

		returnObject["token"] = token
		returnObject["user"] = user

		return c.JSON(returnObject)
	}
}
