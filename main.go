// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"github.com/LunaTearR/go_react_todo/database"
	"github.com/LunaTearR/go_react_todo/handlers"
)

func main() {
	dbUser := "admin"
	dbName := "todolist"
	dbPassword := "password"
	dbHost := "postgres"

	db, err := database.Connect(dbUser, dbName, dbPassword, dbHost)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = database.SetupTables(db)
	if err != nil {
		log.Fatalln("Error setting up tables:", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders: []string{"Origin, Content-Type, Accept"},
	}))

	app.Get("/healthcheck", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Setup Todo routes
	app.Post("/api/todos", func(c fiber.Ctx) error {
		return handlers.CreateTodoHandler(c, db)
	})

	app.Patch("/api/todos/:id/done", func(c fiber.Ctx) error {
		return handlers.MarkTodoDoneHandler(c, db)
	})

	app.Patch("/api/todos/:id/undone", func(c fiber.Ctx) error {
		return handlers.MarkTodoUndoneHandler(c, db)
	})

	app.Get("/api/todos", func(c fiber.Ctx) error {
		return handlers.GetAllTodosHandler(c, db)
	})

	app.Delete("/api/todos/:id/delete", func(c fiber.Ctx) error {
		return handlers.DeleteTodoHandler(c, db)
	})

	app.Get("/users", func(c fiber.Ctx) error {
		return handlers.GetAllUserHandler(c, db)
	})

	app.Post("/users", func(c fiber.Ctx) error {
		return handlers.CreateUserHandler(c, db)
	})

	app.Delete("/users/:id/delete", func(c fiber.Ctx) error {
		return handlers.DeleteUserHandler(c, db)
	})

	log.Fatal(app.Listen(":4000"))
}
