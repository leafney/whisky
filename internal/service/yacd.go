/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-07 17:05
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

func YacdInfo() (*vmodel.Clash, error) {

	// TODO 从配置文件中获取 yacd 端口，如果为空则使用默认值
	port := ""

	fPath, err := utils.LoadByteBashFile(cmds.ScriptYacdStats)
	if err != nil {
		global.GXLog.Errorf("读取 shell 脚本文件失败 [%v]", err)
		return nil, err
	}

	res, err := utils.RunBashFile(fPath, port)
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
		mode := data.Get("mode").String()
		clashInfo.Mode = mode
		clashInfo.AllowLan = data.Get("allow-lan").Bool()

		//	记录下当前最新的 mode 状态
		if err := global.GLevelDB.SetS(vars.KFCYacdMode, rose.StrToLower(mode)); err != nil {
			global.GXLog.Errorf("KFCYacdMode set error [%v]", err)
		}
	} else {
		err2 = fmt.Errorf("获取 clash 信息失败")
	}

	return clashInfo, err2
}

func YacdClashMode(mode string) error {

	// TODO 从配置文件中获取 yacd 端口，如果为空则使用默认值
	port := ""

	fPath, err := utils.LoadByteBashFile(cmds.ScriptYacdMode)
	if err != nil {
		global.GXLog.Errorf("读取 shell 脚本文件失败 [%v]", err)
		return err
	}

	switch mode {
	case vars.ClashModeRule:
		res, err := utils.RunBashFile(fPath, vars.ClashModeRule, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("rule 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	case vars.ClashModeDirect:
		res, err := utils.RunBashFile(fPath, vars.ClashModeDirect, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("direct 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	case vars.ClashModeGlobal:
		res, err := utils.RunBashFile(fPath, vars.ClashModeGlobal, port)
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

	// mode 切换成功，记录下最新的 mode 值，用于自动切换
	if err := global.GLevelDB.SetS(vars.KFCYacdMode, mode); err != nil {
		global.GXLog.Errorf("KFCYacdMode set error [%v]", err)
	}

	return nil
}

func YacdClashSwitch(swt string) error {

	// TODO 从配置文件中获取 yacd 端口，如果为空则使用默认值
	port := ""

	fPath, err := utils.LoadByteBashFile(cmds.ScriptYacdMode)
	if err != nil {
		global.GXLog.Errorf("读取 shell 脚本文件失败 [%v]", err)
		return err
	}

	switch swt {
	case vars.ClashSwitchRule:
		res, err := utils.RunBashFile(fPath, vars.ClashModeRule, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("rule 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	case vars.ClashSwitchDirect:
		res, err := utils.RunBashFile(fPath, vars.ClashModeDirect, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("direct 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}
	default:
		//	自动切换，根据上次的 mode 值，自动判断当前的状态，且只在 rule 和 direct 之间切换；如果是 global 则当做 rule 处理

		// 从缓存中获取上一次的mode状态
		lastMode, _ := global.GLevelDB.GetS(vars.KFCYacdMode)
		nextMode := ""
		if rose.StrEqualFold(lastMode, vars.ClashModeDirect) {
			// 如果上一次为 direct，则新状态为 rule
			nextMode = vars.ClashModeRule
		} else {
			// 如果上一次非 direct(空、rule、global)，则都认为是 rule,新状态为 direct
			nextMode = vars.ClashModeDirect
		}

		res, err := utils.RunBashFile(fPath, nextMode, port)
		if err != nil {
			return err
		}

		global.GXLog.Infof("direct 返回的状态码 [%v]", res)
		if !rose.StrEqualFold(res, vars.ClashStatusCode) {
			return fmt.Errorf("操作异常，返回的状态码为 [%v]", res)
		}

		// 操作成功，记录下最新的 mode 值
		if err := global.GLevelDB.SetS(vars.KFCYacdMode, nextMode); err != nil {
			global.GXLog.Errorf("KFCYacdMode set error [%v]", err)
		}
	}

	return nil
}
