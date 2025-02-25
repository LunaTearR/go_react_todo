package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders: []string{"Origin, Content-Type, Accept"},
	}))

	todos := []Todo{}

	app.Get("/healthcheck", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/api/todos", func(c fiber.Ctx) error {
		todo := &Todo{}

		if err := c.Bind().JSON(todo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.JSON(todo)
	})

	app.Patch("/api/todos/:id/done", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}

		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				return c.JSON(todos[i])
			}
		}

		return c.Status(fiber.StatusNotFound).SendString("Todo not found")
	})

	app.Patch("/api/todos/:id/undone", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}

		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = false
				return c.JSON(todos[i])
			}
		}

		return c.Status(fiber.StatusNotFound).SendString("Todo not found")
	})

	app.Get("/api/todos", func(c fiber.Ctx) error {
		return c.JSON(todos)
	})

	app.Delete("/api/todos/:id/delete", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}

		for i, t := range todos {
			if t.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.JSON(todos)
			}
		}
		return c.Status(fiber.StatusNotFound).SendString("Todo not found")
	})

	log.Fatal(app.Listen(":4000"))
}
