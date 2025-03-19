// handlers/user_handler.go
package handlers

import (
	"strconv"

	"github.com/LunaTearR/go_react_todo/models"
	"github.com/LunaTearR/go_react_todo/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
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

func GetAllUserHandler(c fiber.Ctx, db *sqlx.DB) error {
	user := []models.User{}

	err := db.Select(&user, "SELECT id, username, email FROM users ORDER BY id")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
	
}

func DeleteUserHandler(c fiber.Ctx, db *sqlx.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	query := `DELETE FROM users WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error:": err.Error()})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}