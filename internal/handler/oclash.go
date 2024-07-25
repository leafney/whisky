/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-07 19:23
 * @Description: openclash 相关操作
 */

package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/global/response"
	"github.com/leafney/whisky/global/vars"
	"github.com/leafney/whisky/internal/service"
)

// TODO 待实现
func OClashAction(c fiber.Ctx) error {
	var data map[string]string
	if err := c.Bind().JSON(&data); err != nil {
		global.GXLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	global.GXLog.Info(data)

	if status, ok := data[vars.ClashStatus]; ok {
		global.GXLog.Infof("status %v", status)
		if err := service.SCrashStatus(status); err != nil {
			global.GXLog.Errorf("SCrashStatus error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		global.GXLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}

func OClashRestart(c fiber.Ctx) error {

	return response.Ok(c)
}
