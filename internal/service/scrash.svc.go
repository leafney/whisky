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

	"github.com/leafney/rose"
	"github.com/leafney/whisky/config/vars"
	"github.com/leafney/whisky/pkg/cmds"
	"github.com/leafney/whisky/pkg/utils"
	"github.com/leafney/whisky/pkg/xlogx"
)

type SCrash struct {
	XLog *xlogx.XLogSvc
}

func (s *SCrash) SCrashStatus(status string) error {
	// 判断是否存在 shellCrash 服务
	if exist := s.ChkCrashExist(); !exist {
		return errors.New("shellCrash not found")
	}

	switch status {
	case vars.ClashStsStart:
	case vars.ClashStsRestart:
		go func() {
			if _, err := utils.RunBash(cmds.ScriptCrashStart); err != nil {
				s.XLog.Errorf("shell 脚本 [ScriptCrashStart] 执行失败 [%v]", err)
			}
		}()
	case vars.ClashStsStop:
		go func() {
			if _, err := utils.RunBash(cmds.ScriptCrashStop); err != nil {
				s.XLog.Errorf("shell 脚本 [ScriptCrashStop] 执行失败 [%v]", err)
			}
		}()
	default:
		return errors.New("不支持的操作状态")
	}

	return nil
}

func (s *SCrash) ClashTest() {
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

	s.XLog.Infof("res [%v] err [%v]", res, err)
}

func (s *SCrash) ChkCrashExist() bool {
	fPath, err := utils.LoadByteBashFile(cmds.ScriptCrashExist)
	if err != nil {
		s.XLog.Errorf("shell 脚本 [ScriptCrashExist] 载入失败 [%v]", err)
		return false
	}
	res, err := utils.RunBashFile(fPath)
	if err != nil {
		s.XLog.Errorf("shell 脚本 [ScriptCrashExist] 执行失败 [%v]", err)
		return false
	}

	s.XLog.Debugf("shell 脚本 [ScriptCrashExist] 执行结果 [%v]", res)

	return rose.StrToBool(res)
}
