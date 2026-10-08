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
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/parsex"
	"github.com/leafney/whisky/pkg/response"
	"github.com/leafney/whisky/pkg/xlogx"
)

type CronTask struct {
	XLog           *xlogx.XLogSvc
	CronSvc        *service.Cron // service 层，用于获取通用任务方法
	NetworkTaskBiz *biz.NetworkTask
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
	case biz.NetworkMonitorTaskID:
		return a.getNetworkMonitorStatus(c)
	default:
		return response.Fail(c, "不支持的任务类型")
	}
}

// ControlTask 控制任务（启动/停止/重启）
func (a *CronTask) ControlTask(c *fiber.Ctx) error {
	taskId := c.Params("taskId")
	if taskId == "" {
		a.XLog.Error("任务控制请求缺少任务ID参数")
		return response.Fail(c, "任务ID不能为空")
	}

	var req struct {
		Action string                       `json:"action"` // start, stop, restart, run_now, reset
		Config *vmodel.NetworkMonitorConfig `json:"config,omitempty"`
	}

	if err := parsex.ParseAll(c, &req); err != nil {
		a.XLog.Errorf("解析任务控制请求参数失败: %v", err)
		return response.Fail(c, "请求参数无效")
	}

	a.XLog.Infof("定时任务控制请求: 任务ID=%s, 操作=%s", taskId, req.Action)

	switch taskId {
	case biz.NetworkMonitorTaskID:
		return a.controlNetworkMonitor(c, req.Action, req.Config)
	default:
		a.XLog.Errorf("不支持的任务类型: %s", taskId)
		return response.Fail(c, "不支持的任务类型: "+taskId)
	}
}

// 获取网络监控状态
func (a *CronTask) getNetworkMonitorStatus(c *fiber.Ctx) error {
	// 检查功能是否启用
	featureStatus := a.NetworkTaskBiz.GetTaskFeatureStatus()
	if !featureStatus["feature_enabled"].(bool) {
		return response.OkWithData(c, map[string]interface{}{
			"feature_enabled": false,
			"message":         "网络监控功能未启用，请在配置文件中启用该功能",
		})
	}

	result, err := a.NetworkTaskBiz.GetTaskStatus()
	if err != nil {
		a.XLog.Errorf("获取网络监控状态失败: %v", err)
		return response.Fail(c, err.Error())
	}

	a.XLog.Debugf("网络监控状态查询完成 - 任务运行: %v", result["task_running"])
	return response.OkWithData(c, result)
}

// 控制网络监控任务
func (a *CronTask) controlNetworkMonitor(c *fiber.Ctx, action string, config *vmodel.NetworkMonitorConfig) error {
	// 检查功能是否启用
	featureStatus := a.NetworkTaskBiz.GetTaskFeatureStatus()
	if !featureStatus["feature_enabled"].(bool) {
		a.XLog.Errorf("网络监控功能未启用，无法执行定时任务操作: %s", action)
		return response.Fail(c, "网络监控功能未启用，请在配置文件中启用该功能后重启服务")
	}

	var err error
	var message string

	switch action {
	case "start":
		err = a.NetworkTaskBiz.StartTask()
		if err != nil {
			a.XLog.Errorf("启动网络监控任务失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络监控任务已启动"
		a.XLog.Info(message)

	case "stop":
		err = a.NetworkTaskBiz.StopTask()
		if err != nil {
			a.XLog.Errorf("停止网络监控任务失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络监控任务已停止"
		a.XLog.Info(message)

	case "restart":
		err = a.NetworkTaskBiz.RestartTask()
		if err != nil {
			a.XLog.Errorf("重启网络监控任务失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络监控任务已重启"
		a.XLog.Info(message)

	case "run_now":
		err = a.NetworkTaskBiz.RunTaskNow()
		if err != nil {
			a.XLog.Errorf("立即执行网络监控失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络监控任务已立即执行"
		a.XLog.Info(message)

	case "reset":
		err = a.NetworkTaskBiz.ResetTask()
		if err != nil {
			a.XLog.Errorf("重置网络监控状态失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络监控状态已重置"
		a.XLog.Info(message)

	default:
		a.XLog.Errorf("无效的定时任务操作: %s", action)
		return response.Fail(c, "无效的操作类型，支持的操作: start, stop, restart, run_now, reset")
	}

	return response.OkWithData(c, map[string]interface{}{
		"task_id": biz.NetworkMonitorTaskID,
		"action":  action,
		"message": message,
	})
}

// GetAllTasksStatus 获取所有任务状态
func (a *CronTask) GetAllTasksStatus(c *fiber.Ctx) error {
	// 获取网络监控的完整状态
	networkMonitorStatus := a.NetworkTaskBiz.GetTaskFeatureStatus()

	tasks := []map[string]interface{}{
		{
			"task_id":         biz.NetworkMonitorTaskID,
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
