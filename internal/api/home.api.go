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
	"github.com/leafney/whisky/pkg/version"
)

type Home struct {
	VersionInfo *version.Info
}

func (a *Home) Home(c *fiber.Ctx) error {
	return response.OkWithData(c, "Hello Whisky!")
}

func (a *Home) Version(c *fiber.Ctx) error {
	return response.OkWithData(c, version.GetVersion())
}
