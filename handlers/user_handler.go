// handlers/user_handler.go
package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/LunaTearR/go_react_todo/models"
	"github.com/LunaTearR/go_react_todo/utils"
)

// CreateUserHandler handles user creation
func CreateUserHandler(c fiber.Ctx, db *sqlx.DB) error {
	user := models.User{}

	// Parse JSON request
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hash password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Insert user into the database
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRowx(query, user.Username, user.Email, hashedPassword).Scan(&user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user.Password = ""

	return c.JSON(user)
}