/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:55
 * @Description:
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/config/vars"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/pkg/response"
	"github.com/leafney/whisky/pkg/xlogx"
)

type OClash struct {
	XLog *xlogx.XLogSvc
}

// TODO 待实现
func (a *OClash) OClashAction(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.Bind().JSON(&data); err != nil {
		a.XLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	a.XLog.Info(data)

	if status, ok := data[vars.ClashStatus]; ok {
		a.XLog.Infof("status %v", status)
		if err := service.SCrashStatus(status); err != nil {
			a.XLog.Errorf("SCrashStatus error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		a.XLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}

func (a *OClash) OClashRestart(c *fiber.Ctx) error {

	return response.Ok(c)
}
