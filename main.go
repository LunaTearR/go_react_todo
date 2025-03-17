package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

type User struct {
	Name  string `db:"username"`
	Email string `db:"email"`
}

func main() {

	dbUser := "admin"
	dbName := "todolist"
	dbPassword := "password"
	dbHost := "postgres"

	// Construct the connection string
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s",
		dbUser, dbName, dbPassword, dbHost)

	log.Printf("Attempting to connect to PostgreSQL: host=%s, user=%s, dbname=%s", dbHost, dbUser, dbName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to database")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT,
			done BOOLEAN DEFAULT FALSE
		)
	`)
	
	if err != nil {
		log.Println("Error creating todos table:", err)
	} else {
		log.Println("Table 'todos' created or already exists")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`)
	
	if err != nil {
		log.Println("Error creating users table:", err)
	} else {
		log.Println("Table 'users' created or already exists")
	}

	place := User{}
	rows, _ := db.Queryx("SELECT username, email FROM users")
	for rows.Next() {
		err := rows.StructScan(&place)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%#v\n", place)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders: []string{"Origin, Content-Type, Accept"},
	}))

	app.Get("/healthcheck", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/api/todos", func(c fiber.Ctx) error {
		todo := &Todo{}

		if err := c.Bind().JSON(todo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		query := `INSERT INTO todos (title, body, done) VALUES ($1, $2, $3) RETURNING id`
		err := db.QueryRowx(query, todo.Title, todo.Body, todo.Done).Scan(&todo.ID)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(todo)
	})

	app.Patch("/api/todos/:id/done", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}

		query := `UPDATE todos SET done = true WHERE id = $1 RETURNING id, title, body, done`
		todo := Todo{}
		err = db.QueryRowx(query, id).StructScan(&todo)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Todo not found")
		}

		return c.JSON(todo)
	})

	app.Patch("/api/todos/:id/undone", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
		}

		query := `UPDATE todos SET done = false WHERE id = $1 RETURNING id, title, body, done`
		todo := Todo{}
		err = db.QueryRowx(query, id).StructScan(&todo)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Todo not found")
		}

		return c.JSON(todo)
	})

	app.Get("/api/todos", func(c fiber.Ctx) error {
		todos := []Todo{}
		err := db.Select(&todos, "SELECT id, title, body, done FROM todos ORDER BY id")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(todos)
	})

	app.Delete("/api/todos/:id/delete", func(c fiber.Ctx) error {
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

		// Return all remaining todos after deletion
		todos := []Todo{}
		err = db.Select(&todos, "SELECT id, title, body, done FROM todos ORDER BY id")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(todos)
	})

	log.Fatal(app.Listen(":4000"))
}
