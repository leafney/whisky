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
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/global/vars"
	"github.com/leafney/whisky/pkgs/cmds"
	"github.com/leafney/whisky/utils"
)

func SCrashStatus(status string) error {
	// TODO 判断是否存在 shellCrash 服务

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
		return errors.New("不支持的操作状态")
	}

	return nil
}

func ClashTest() {
	//fPath, err := utils.LoadByteBashFile(cmds.ScriptYacdModeB)
	//if err != nil {
	//	global.GXLog.Errorf("读取 shell 脚本文件失败 [%v]", err)
	//	return
	//}
	//global.GXLog.Infof("shell 脚本文件 [%v]", fPath)
	//
	//res, err := utils.RunBashFile(fPath, "direct")

	//res, err := utils.RunBash("hello='$1'; echo $hello", "hello")

	command := "hello='nihao'; echo -n $hello"

	// 要传入的参数
	//param := "world"

	//// 创建命令
	//cmd := exec.Command("/bin/sh", "-c", command)
	//
	//// 设置命令参数
	//cmd.Args = append(cmd.Args, param)
	//
	//// 执行命令并获取输出
	//res, err := cmd.CombinedOutput()

	//command = fmt.Sprintf(command, param)
	res, err := utils.RunBash(command)

	global.GXLog.Infof("res [%v] err [%v]", res, err)
}
