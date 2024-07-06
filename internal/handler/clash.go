/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 18:38
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

func ClashInfo(c fiber.Ctx) error {

	info, err := service.ClashInfo()
	if err != nil {
		return response.Fail(c, err.Error())
	}

	return response.OkWithData(c, info)
}

func ClashAction(c fiber.Ctx) error {

	var data map[string]string
	if err := c.Bind().JSON(&data); err != nil {
		global.GXLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	global.GXLog.Info(data)

	if status, ok := data[vars.ClashStatus]; ok {
		global.GXLog.Infof("status %v", status)
		if err := service.ClashStatus(status); err != nil {
			global.GXLog.Errorf("ClashStatus error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else if mode, ok := data[vars.ClashMode]; ok {
		global.GXLog.Infof("mode %v", mode)
		if err := service.ClashMode(mode); err != nil {
			global.GXLog.Errorf("ClashMode error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else if swt, ok := data[vars.ClashSwitch]; ok {
		global.GXLog.Infof("switch %v", swt)
		if err := service.ClashSwitch(swt); err != nil {
			global.GXLog.Errorf("ClashSwitch error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		global.GXLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}
