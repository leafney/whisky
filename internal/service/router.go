/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 23:08
 * @Description:
 */

package service

import (
	"github.com/leafney/rose"
	"github.com/leafney/whisky/global"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/cmds"
	"github.com/leafney/whisky/pkg/utils"
)

func RouterInfo() *vmodel.Stat {

	statInfo := new(vmodel.Stat)

	//	获取 cpu 温度
	cpuTemp, err := utils.RunBash(cmds.ScriptTempCpu)
	if err != nil {
		global.GXLog.Errorf("shell 脚本 [ScriptTempCpu] 执行失败 [%v]", err)
		cpuTemp = ""
	}
	statInfo.CpuTemp = cpuTemp

	// mem usage
	memUsage, err := utils.RunBash(cmds.ScriptMemUsage)
	if err != nil {
		global.GXLog.Errorf("shell 脚本 [ScriptMemUsage] 执行失败 [%v]", err)
		memUsage = ""
	}
	statInfo.MemUsage = memUsage

	//	disk usage
	diskUsage, err := utils.RunBash(cmds.ScriptDiskUsage)
	if err != nil {
		global.GXLog.Errorf("shell 脚本 [ScriptDiskUsage] 执行失败 [%v]", err)
		diskUsage = ""
	}
	statInfo.DiskUsage = diskUsage

	//	running time
	runTime, err := utils.RunBash(cmds.ScriptRunningTime)
	if err != nil {
		global.GXLog.Errorf("shell 脚本 [ScriptRunningTime] 执行失败 [%v]", err)
		runTime = ""
	}
	statInfo.RunningTime = runTime

	//	boot time
	bootTime, err := utils.RunBash(cmds.ScriptBootTime)
	if err != nil {
		global.GXLog.Errorf("shell 脚本 [ScriptBootTime] 执行失败 [%v]", err)
		bootTime = ""
	}
	statInfo.BootTime = bootTime

	//	now time
	nowTime, err := utils.RunBash(cmds.ScriptTimeNow)
	if err != nil {
		global.GXLog.Errorf("shell 脚本 [ScriptTimeNow] 执行失败 [%v]", err)
		nowTime = rose.TNowDateTime()
	}
	statInfo.NowTime = nowTime

	return statInfo
}

func RouterRestart() error {

	go func() {
		if _, err := utils.RunBash(cmds.ScriptReboot); err != nil {
			global.GXLog.Errorf("shell 脚本 [ScriptReboot] 执行失败 [%v]", err)
		}
	}()

	return nil
}
