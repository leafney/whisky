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
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/global/response"
	"github.com/leafney/whisky/global/vars"
	"github.com/leafney/whisky/internal/service"
)

func RouterInfo(c fiber.Ctx) error {
	stat := service.RouterInfo()
	return response.OkWithData(c, stat)
}

func RouterStatus(c fiber.Ctx) error {
	var data map[string]string
	if err := c.Bind().JSON(&data); err != nil {
		global.GXLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	global.GXLog.Info(data)

	if status, ok := data[vars.RouterStatus]; ok && status == vars.RouterStsRestart {
		global.GXLog.Infof("status %v", status)
		if err := service.RouterRestart(); err != nil {
			global.GXLog.Errorf("RouterStatus error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		global.GXLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}
