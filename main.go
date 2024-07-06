/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-05 23:31
 * @Description:
 */

package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leafney/whisky/cmd/core"
	"github.com/leafney/whisky/internal/handler"
)

func main() {

	core.InitXLog()

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello Whisky!")
	})

	app.Get("/cpu", handler.GetCpuTemp)

	app.Listen(":8080")
}
