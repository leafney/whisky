/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-07 19:18
 * @Description:
 */

package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/global/response"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello Whisky!")
}

func Version(c *fiber.Ctx) error {

	return response.Ok(c)
}
