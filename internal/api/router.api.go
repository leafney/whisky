/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:56
 * @Description:
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/config/vars"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/pkg/parsex"
	"github.com/leafney/whisky/pkg/response"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Router struct {
	XLog      *xlogx.XLogSvc
	RouterSvc *service.Router
}

func (a *Router) RouterInfo(c *fiber.Ctx) error {
	stat := a.RouterSvc.RouterInfo()
	return response.OkWithData(c, stat)
}

func (a *Router) RouterStatus(c *fiber.Ctx) error {
	var data map[string]string
	if err := parsex.ParseAll(c, &data); err != nil {
		a.XLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	a.XLog.Info(data)

	if status, ok := data[vars.RouterStatus]; ok && status == vars.RouterStsRestart {
		a.XLog.Infof("status %v", status)
		if err := a.RouterSvc.RouterRestart(); err != nil {
			a.XLog.Errorf("RouterStatus error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		a.XLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}
