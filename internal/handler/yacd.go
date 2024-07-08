/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-07 17:10
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

func YacdClashInfo(c fiber.Ctx) error {

	info, err := service.YacdInfo()
	if err != nil {
		return response.Fail(c, err.Error())
	}

	return response.OkWithData(c, info)
}

func YacdClashAction(c fiber.Ctx) error {

	var data map[string]string
	if err := c.Bind().JSON(&data); err != nil {
		global.GXLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	global.GXLog.Info(data)

	if mode, ok := data[vars.ClashMode]; ok {
		global.GXLog.Infof("mode %v", mode)
		if err := service.YacdClashMode(mode); err != nil {
			global.GXLog.Errorf("ClashMode error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else if swt, ok := data[vars.ClashSwitch]; ok {
		global.GXLog.Infof("switch %v", swt)
		if err := service.YacdClashSwitch(swt); err != nil {
			global.GXLog.Errorf("ClashSwitch error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else if lan, ok := data[vars.ClashLan]; ok {
		global.GXLog.Infof("lan %v", lan)
		if err := service.YacdClashAllowLan(lan); err != nil {
			global.GXLog.Errorf("ClashLan error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		global.GXLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}
