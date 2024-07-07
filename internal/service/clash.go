/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 19:03
 * @Description:
 */

package service

import (
	"errors"
	"fmt"
	"github.com/leafney/rose"
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/global/vars"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkgs/cmds"
	"github.com/leafney/whisky/utils"
	"github.com/tidwall/gjson"
)

func ClashInfo() (*vmodel.Clash, error) {

	res, err := utils.RunBash(cmds.ScriptYacdStats)
	if err != nil {
		return nil, err
	}

	//global.GXLog.Infof("clash info [%v]", res)

	clashInfo := new(vmodel.Clash)
	var err2 error
	if !rose.StrIsEmpty(res) {
		data := gjson.Parse(res)
		clashInfo.HttpPort = data.Get("port").String()
		clashInfo.SocksPort = data.Get("socks-port").String()
		clashInfo.MixedPort = data.Get("mixed-port").String()
		clashInfo.RedirPort = data.Get("redir-port").String()
		clashInfo.Mode = data.Get("mode").String()
		clashInfo.AllowLan = data.Get("allow-lan").Bool()
	} else {
		err2 = fmt.Errorf("获取 clash 信息失败")
	}

	return clashInfo, err2
}

func ClashStatus(status string) error {
	// TODO 判断是否存在 crash 服务

	switch status {
	case vars.ClashStsStart:
	case vars.ClashStsRestart:
		_, err := utils.RunBash(cmds.ScriptCrashStart)
		if err != nil {
			return err
		}
	case vars.ClashStsStop:
		if _, err := utils.RunBash(cmds.ScriptCrashStop); err != nil {
			return err
		}
	default:
		return errors.New("不支持的状态")
	}

	return nil
}

func ClashMode(mode string) error {
	// TODO 判断是否存在 crash 服务

	// TODO 从配置文件中获取 yacd 端口，如果为空则使用默认值
	port := ""

	switch mode {
	case vars.ClashModeRule:
		res, err := utils.RunBash(cmds.ScriptYacdMode, vars.ClashModeRule, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("rule 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	case vars.ClashModeDirect:
		res, err := utils.RunBash(cmds.ScriptYacdMode, vars.ClashModeDirect, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("direct 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	case vars.ClashModeGlobal:
		res, err := utils.RunBash(cmds.ScriptYacdMode, vars.ClashModeGlobal, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("global 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	default:
		return errors.New("不支持的 Mode")
	}

	// TODO mode 切换成功，记录下最新的 mode 值，用于自动切换

	return nil
}

func ClashSwitch(swt string) error {
	// TODO 判断是否存在 crash 服务

	// TODO 从配置文件中获取 yacd 端口，如果为空则使用默认值
	port := ""

	switch swt {
	case vars.ClashSwitchRule:
		res, err := utils.RunBash(cmds.ScriptYacdMode, vars.ClashModeRule, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("rule 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	case vars.ClashSwitchDirect:
		res, err := utils.RunBash(cmds.ScriptYacdMode, vars.ClashModeDirect, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("direct 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	default:
		//	自动切换，根据上次的 mode 值，自动判断当前的状态，且只在 rule 和 direct 之间切换；如果是 global 则当做 rule 处理

		//	TODO 从缓存中获取上一次的mode状态
		lastMode := ""
		nextMode := ""
		if rose.StrEqualFold(lastMode, vars.ClashModeDirect) {
			// 如果上一次为 direct，则新状态为 rule
			nextMode = vars.ClashModeRule
		} else {
			// 如果上一次非 direct(空、rule、global)，则都认为是 rule,新状态为 direct
			nextMode = vars.ClashModeDirect
		}

		res, err := utils.RunBash(cmds.ScriptYacdMode, nextMode, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("direct 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	}

	return nil
}

func ClashTest() {
	fPath, err := utils.LoadByteBashFile(cmds.ScriptYacdModeB)
	if err != nil {
		global.GXLog.Errorf("读取 shell 脚本文件失败 [%v]", err)
		return
	}
	global.GXLog.Infof("shell 脚本文件 [%v]", fPath)

	res, err := utils.RunBashFile(fPath, "direct")

	global.GXLog.Infof("res [%v] err [%v]", res, err)
}
