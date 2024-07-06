/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:37
 * @Description:
 */

package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leafney/whisky/internal/service"
)

func GetCpuTemp(c fiber.Ctx) error {

	temp := service.GetCpuTemp()

	return c.SendString(temp)
}
