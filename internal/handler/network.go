/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 18:37
 * @Description:
 */

package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leafney/whisky/global/response"
	"github.com/leafney/whisky/internal/service"
)

func NetWorkInfo(c fiber.Ctx) error {
	service.ClashTest()
	return response.Ok(c)
}
