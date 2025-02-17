/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:48
 * @Description:
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/pkg/response"
)

type NetWork struct {
}

func (a *NetWork) NetWorkInfo(c *fiber.Ctx) error {
	netInfo, err := service.NetWorkInfo()
	if err != nil {
		return response.Fail(c, err.Error())
	}
	return response.OkWithData(c, netInfo)
}
