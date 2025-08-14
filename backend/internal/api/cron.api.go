/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 定时任务控制API层
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/parsex"
	"github.com/leafney/whisky/pkg/response"
	"github.com/leafney/whisky/pkg/xlogx"
)

type CronTask struct {
	XLog              *xlogx.XLogSvc
	CronSvc           *service.Cron
	NetworkMonitorSvc *service.NetworkMonitor
}

// GetTaskMethods 获取所有可用的任务方法
func (a *CronTask) GetTaskMethods(c *fiber.Ctx) error {
	methods := a.CronSvc.GetAvailableMethods()
	return response.OkWithData(c, map[string]interface{}{
		"methods": methods,
		"count":   len(methods),
	})
}

// GetTaskStatus 获取任务状态
func (a *CronTask) GetTaskStatus(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	if taskId == "" {
		return response.Fail(c, "任务ID不能为空")
	}

	switch taskId {
	case service.NetworkMonitorTaskID:
		return a.getNetworkMonitorStatus(c)
	default:
		return response.Fail(c, "不支持的任务类型")
	}
}

// ControlTask 控制任务（启动/停止/重启）
func (a *CronTask) ControlTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	if taskId == "" {
		return response.Fail(c, "任务ID不能为空")
	}

	var req struct {
		Action string                       `json:"action"` // start, stop, restart, run_now, update_config
		Config *vmodel.NetworkMonitorConfig `json:"config,omitempty"`
	}

	if err := parsex.ParseAll(c, &req); err != nil {
		a.XLog.Errorf("解析任务控制请求参数失败: %v", err)
		return response.Fail(c, "请求参数无效")
	}

	a.XLog.Infof("任务控制请求: 任务=%s, 操作=%s", taskId, req.Action)

	switch taskId {
	case service.NetworkMonitorTaskID:
		return a.controlNetworkMonitor(c, req.Action, req.Config)
	default:
		return response.Fail(c, "不支持的任务类型")
	}
}

// 获取网络监控状态
func (a *CronTask) getNetworkMonitorStatus(c *fiber.Ctx) error {
	// 获取任务运行状态
	isRunning := a.CronSvc.IsNetworkMonitorRunning()

	// 获取监控详细状态（通过 NetworkMonitorSvc）
	status, err := a.NetworkMonitorSvc.GetMonitorStatus()
	if err != nil {
		a.XLog.Errorf("获取网络监控状态失败: %v", err)
		return response.Fail(c, "获取监控状态失败")
	}

	stats, err := a.NetworkMonitorSvc.GetMonitorStats()
	if err != nil {
		a.XLog.Errorf("获取网络监控统计失败: %v", err)
		return response.Fail(c, "获取监控统计失败")
	}

	result := map[string]interface{}{
		"task_id":    service.NetworkMonitorTaskID,
		"is_running": isRunning,
		"status":     status,
		"stats":      stats,
	}

	return response.OkWithData(c, result)
}

// 控制网络监控任务
func (a *CronTask) controlNetworkMonitor(c *fiber.Ctx, action string, config *vmodel.NetworkMonitorConfig) error {
	switch action {
	case "start":
		if err := a.CronSvc.StartNetworkMonitor(); err != nil {
			a.XLog.Errorf("启动网络监控任务失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控任务已启动")

	case "stop":
		if err := a.CronSvc.StopNetworkMonitor(); err != nil {
			a.XLog.Errorf("停止网络监控任务失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控任务已停止")

	case "restart":
		if err := a.CronSvc.RestartNetworkMonitor(); err != nil {
			a.XLog.Errorf("重启网络监控任务失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控任务已重启")

	case "run_now":
		if err := a.CronSvc.RunNetworkMonitorNow(); err != nil {
			a.XLog.Errorf("立即执行网络监控失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控任务已立即执行")

	case "reset":
		if err := a.NetworkMonitorSvc.ResetMonitor(); err != nil {
			a.XLog.Errorf("重置网络监控状态失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控状态已重置")

	default:
		a.XLog.Errorf("无效的任务操作: %s", action)
		return response.Fail(c, "无效的操作类型")
	}

	return response.Ok(c)
}

// GetAllTasksStatus 获取所有任务状态
func (a *CronTask) GetAllTasksStatus(c *fiber.Ctx) error {
	// 获取网络监控的完整状态
	networkMonitorStatus := a.CronSvc.GetNetworkMonitorFeatureStatus()

	tasks := []map[string]interface{}{
		{
			"task_id":         service.NetworkMonitorTaskID,
			"name":            "网络连通性监控",
			"is_running":      networkMonitorStatus["task_running"],
			"feature_enabled": networkMonitorStatus["feature_enabled"],
			"config_valid":    networkMonitorStatus["config_valid"],
		},
	}

	return response.OkWithData(c, map[string]interface{}{
		"tasks": tasks,
		"count": len(tasks),
	})
}
