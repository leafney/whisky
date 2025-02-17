/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-07 12:11
 * @Description:
 */

package core

import (
	"github.com/leafney/whisky/config/vars"
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/pkg/utils"
)

func InitShellClean() {
	//	在程序启动时清空 `.shell` 目录

	if err := utils.DeleteFilesByExtension(vars.ShellTempDir, ".sky"); err != nil {
		global.GXLog.Errorf("清空 shell 脚本临时目录操作异常 [%v]", err)
	}

	global.GXLog.Infoln("[ShellClean] Load successful")
}
