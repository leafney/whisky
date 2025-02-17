/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:17
 * @Description:
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/pkg/response"
)

type Home struct {
}

func (a *Home) Home(c *fiber.Ctx) error {
	return response.OkWithData(c, "Hello Whisky!")
}

func (a *Home) Version(c *fiber.Ctx) error {
	return response.Ok(c)
}
