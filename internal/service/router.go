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
	"github.com/leafney/whisky/pkgs/cmds"
	"github.com/leafney/whisky/utils"
)

func RouterInfo() *vmodel.Stat {

	statInfo := new(vmodel.Stat)

	//	获取 cpu 温度
	cpuTemp, err := utils.RunBash(cmds.ScriptTempCpu)
	if err != nil {
		cpuTemp = ""
	}
	statInfo.CpuTemp = cpuTemp

	// mem usage
	memUsage, err := utils.RunBash(cmds.ScriptMemUsage)
	if err != nil {
		memUsage = ""
	}
	statInfo.MemUsage = memUsage

	//	disk usage
	diskUsage, err := utils.RunBash(cmds.ScriptDiskUsage)
	if err != nil {
		diskUsage = ""
	}
	statInfo.DiskUsage = diskUsage

	//	running time
	runTime, err := utils.RunBash(cmds.ScriptRunningTime)
	if err != nil {
		runTime = ""
	}
	statInfo.RunningTime = runTime

	//	boot time
	bootTime, err := utils.RunBash(cmds.ScriptBootTime)
	if err != nil {
		bootTime = ""
	}
	statInfo.BootTime = bootTime

	//	now time
	nowTime, err := utils.RunBash(cmds.ScriptTimeNow)
	if err != nil {
		nowTime = rose.TNowDateTime()
	}
	statInfo.NowTime = nowTime

	return statInfo
}

func RouterRestart() error {

	go func() {
		if _, err := utils.RunBash(cmds.ScriptReboot); err != nil {
			global.GXLog.Errorf("ScriptReboot error [%v]", err)
		}
	}()

	return nil
}
