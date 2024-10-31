package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/hello", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello from the serviceOne")
	})

	app.Post("/hello/:name", func(ctx *fiber.Ctx) error {
		name := ctx.Query("name")
		return ctx.SendString(fmt.Sprintf("Hello: %s from serviceOne", name))
	})

	err := app.Listen(":3001")
	if err != nil {
		return
	}
}
