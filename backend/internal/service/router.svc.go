/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 23:08
 * @Description:
 */

package service

import (
	"context"
	"strings"
	"time"

	"github.com/leafney/rose"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/cmds"
	"github.com/leafney/whisky/pkg/utils"
	"github.com/leafney/whisky/pkg/xlogx"
	"golang.org/x/sync/errgroup"
)

type Router struct {
	XLog *xlogx.XLogSvc
}

func (s *Router) RouterInfo() *vmodel.Stat {
	// 并发 + 超时 收集信息
	timeout := 3 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	statInfo := new(vmodel.Stat)

	var cpuTemp, memUsage, diskUsage, runTime, bootTime, nowTime string

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		out, err := utils.RunBashCtx(ctx, cmds.ScriptTempCpu)
		if err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptTempCpu] 执行失败 [%v]", err)
			cpuTemp = ""
			return nil
		}
		cpuTemp = strings.TrimSpace(out)
		return nil
	})
	g.Go(func() error {
		out, err := utils.RunBashCtx(ctx, cmds.ScriptMemUsage)
		if err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptMemUsage] 执行失败 [%v]", err)
			memUsage = ""
			return nil
		}
		memUsage = strings.TrimSpace(out)
		return nil
	})
	g.Go(func() error {
		out, err := utils.RunBashCtx(ctx, cmds.ScriptDiskUsage)
		if err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptDiskUsage] 执行失败 [%v]", err)
			diskUsage = ""
			return nil
		}
		diskUsage = strings.TrimSpace(out)
		return nil
	})
	g.Go(func() error {
		out, err := utils.RunBashCtx(ctx, cmds.ScriptRunningTime)
		if err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptRunningTime] 执行失败 [%v]", err)
			runTime = ""
			return nil
		}
		runTime = strings.TrimSpace(out)
		return nil
	})
	g.Go(func() error {
		out, err := utils.RunBashCtx(ctx, cmds.ScriptBootTime)
		if err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptBootTime] 执行失败 [%v]", err)
			bootTime = ""
			return nil
		}
		bootTime = strings.TrimSpace(out)
		return nil
	})
	g.Go(func() error {
		out, err := utils.RunBashCtx(ctx, cmds.ScriptTimeNow)
		if err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptTimeNow] 执行失败 [%v]", err)
			nowTime = rose.TNowDateTime()
			return nil
		}
		nowTime = strings.TrimSpace(out)
		return nil
	})

	_ = g.Wait()

	statInfo.CpuTemp = cpuTemp
	statInfo.MemUsage = memUsage
	statInfo.DiskUsage = diskUsage
	statInfo.RunningTime = runTime
	statInfo.BootTime = bootTime
	statInfo.NowTime = nowTime

	return statInfo
}

func (s *Router) RouterRestart() error {
	go func() {
		// 轻微延时，确保响应能先返回
		delay := 2 * time.Second
		time.Sleep(delay)
		if _, err := utils.RunBash(cmds.ScriptReboot); err != nil {
			s.XLog.Errorf("shell 脚本 [ScriptReboot] 执行失败 [%v]", err)
		}
	}()
	return nil
}
