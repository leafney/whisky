/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 12:40
 * @Description:
 */

package run

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leafney/whisky/internal/handler"
)

func bindRouter(app *fiber.App) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello Whisky!")
	})

	app.Get("/cpu", handler.GetCpuTemp)
}
