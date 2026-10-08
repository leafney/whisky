/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 18:02
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

type YAcd struct {
	XLog    *xlogx.XLogSvc
	YacdSvc *service.YAcd
}

func (a *YAcd) YacdClashInfo(c *fiber.Ctx) error {

	info, err := a.YacdSvc.YacdInfo()
	if err != nil {
		return response.Fail(c, err.Error())
	}

	return response.OkWithData(c, info)
}

func (a *YAcd) YacdClashAction(c *fiber.Ctx) error {

	var data map[string]string
	if err := parsex.ParseJson(c, &data); err != nil {
		a.XLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	a.XLog.Info(data)

	if mode, ok := data[vars.ClashMode]; ok {
		a.XLog.Infof("mode %v", mode)
		if err := a.YacdSvc.YacdClashMode(mode); err != nil {
			a.XLog.Errorf("ClashMode error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else if swt, ok := data[vars.ClashSwitch]; ok {
		a.XLog.Infof("switch %v", swt)
		if err := a.YacdSvc.YacdClashSwitch(swt); err != nil {
			a.XLog.Errorf("ClashSwitch error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else if lan, ok := data[vars.ClashLan]; ok {
		a.XLog.Infof("lan %v", lan)
		if err := a.YacdSvc.YacdClashAllowLan(lan); err != nil {
			a.XLog.Errorf("ClashLan error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		a.XLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}
