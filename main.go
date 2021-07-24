package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Todo struct {
	ID        int    `json:"ID"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{ID: 1, Name: "Walk the dog", Completed: false},
	{ID: 2, Name: "Walk the cat", Completed: true},
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world")
	})

	setupTodoRoutes(app)

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}

}

func setupTodoRoutes(app *fiber.App) {
	todosRoute := app.Group("/todos")

	todosRoute.Get("/", getTodos)
	todosRoute.Get("/:id", getTodo)
	todosRoute.Post("/", createTodo)
	todosRoute.Delete("/:id", deleteTodo)
	todosRoute.Patch("/:id", updateTodo)
}

func updateTodo(c *fiber.Ctx) error {
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := c.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse ID",
		})
	}

	var body request
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
	}

	var todo *Todo

	for _, t := range todos {
		if t.ID == id {
			todo = t
			break
		}
	}

	if todo == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	return c.Status(fiber.StatusOK).JSON(todo)
}

func deleteTodo(c *fiber.Ctx) error {
	paramsId := c.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse ID",
		})
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[0:i], todos[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)
}

func getTodo(c *fiber.Ctx) error {
	paramsId := c.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse ID",
		})
	}

	for _, todo := range todos {
		if todo.ID == id {
			return c.Status(fiber.StatusOK).JSON(todo)
		}
	}

	return c.SendStatus(fiber.StatusNotFound)

}

func createTodo(c *fiber.Ctx) error {
	type request struct {
		Name string `json:"name"`
	}

	var body request
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	todo := &Todo{
		ID:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)
	return c.Status(fiber.StatusCreated).JSON(todo)
}

func getTodos(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}
