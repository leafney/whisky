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

	// router
	app.Get("/router", handler.RouterInfo)
	app.Post("/router", handler.RouterStatus)

	// network
	app.Get("/network", handler.NetWorkInfo)

	//	clash
	app.Get("/clash", handler.ClashInfo)
	app.Post("/clash", handler.ClashAction)

}
