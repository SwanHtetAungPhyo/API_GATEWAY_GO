package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/hello", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello from the serviceTwo")
	})

	app.Post("/hello/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name")
		return ctx.SendString(fmt.Sprintf("Hello: %s from serviceTwo", name))
	})

	err := app.Listen(":3002")
	if err != nil {
		return
	}
}
