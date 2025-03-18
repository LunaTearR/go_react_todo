// handlers/todo_handler.go
package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/LunaTearR/go_react_todo/models"
	"strconv"
)

// CreateTodoHandler handles creation of a new todo
func CreateTodoHandler(c fiber.Ctx, db *sqlx.DB) error {
	todo := &models.Todo{}

	if err := c.Bind().JSON(todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	
	query := `INSERT INTO todos (title, body, done) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRowx(query, todo.Title, todo.Body, todo.Done).Scan(&todo.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(todo)
}

// MarkTodoDoneHandler marks a todo as done
func MarkTodoDoneHandler(c fiber.Ctx, db *sqlx.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	query := `UPDATE todos SET done = true WHERE id = $1 RETURNING id, title, body, done`
	todo := models.Todo{}
	err = db.QueryRowx(query, id).StructScan(&todo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Todo not found")
	}

	return c.JSON(todo)
}

// MarkTodoUndoneHandler marks a todo as not done
func MarkTodoUndoneHandler(c fiber.Ctx, db *sqlx.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	query := `UPDATE todos SET done = false WHERE id = $1 RETURNING id, title, body, done`
	todo := models.Todo{}
	err = db.QueryRowx(query, id).StructScan(&todo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Todo not found")
	}

	return c.JSON(todo)
}

// GetAllTodosHandler retrieves all todos
func GetAllTodosHandler(c fiber.Ctx, db *sqlx.DB) error {
	todos := []models.Todo{}
	err := db.Select(&todos, "SELECT id, title, body, done FROM todos ORDER BY id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(todos)
}

// DeleteTodoHandler deletes a todo
func DeleteTodoHandler(c fiber.Ctx, db *sqlx.DB) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	query := `DELETE FROM todos WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Todo not found")
	}

	todos := []models.Todo{}
	err = db.Select(&todos, "SELECT id, title, body, done FROM todos ORDER BY id")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(todos)
}