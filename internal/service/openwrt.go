/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:30
 * @Description:
 */

package service

import (
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/pkgs/commands"
	"github.com/leafney/whisky/utils"
)

func GetCpuTemp() string {
	res, err := utils.RunBash(commands.ScriptTempCpu)
	if err != nil {
		global.GXLog.Errorf("获取 cpu 温度操作失败 [%v]", err)
		return ""
	}

	global.GXLog.Infof("当前 cpu 温度 [%v]", res)
	return res
}
