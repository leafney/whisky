/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:56
 * @Description:
 */

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/config/vars"
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/parsex"
	"github.com/leafney/whisky/pkg/response"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Router struct {
	XLog           *xlogx.XLogSvc
	RouterSvc      *service.Router
	NetworkTaskBiz *biz.NetworkTask
}

func (a *Router) RouterInfo(c *fiber.Ctx) error {
	stat := a.RouterSvc.RouterInfo()
	return response.OkWithData(c, stat)
}

func (a *Router) RouterStatus(c *fiber.Ctx) error {
	var data map[string]string
	if err := parsex.ParseAll(c, &data); err != nil {
		a.XLog.Errorf("解析路由器状态请求参数异常: %v", err)
		return response.Fail(c, "Invalid request body")
	}

	a.XLog.Debugf("路由器状态控制请求: %+v", data)

	if status, ok := data[vars.RouterStatus]; ok && status == vars.RouterStsRestart {
		a.XLog.Infof("执行路由器重启操作: %s", status)
		if err := a.RouterSvc.RouterRestart(); err != nil {
			a.XLog.Errorf("路由器重启操作失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("路由器重启命令已执行")
	} else {
		a.XLog.Errorf("路由器状态控制参数错误: %+v", data)
		return response.Fail(c, "参数错误，需要正确的路由器状态控制参数")
	}

	return response.Ok(c)
}

// NetworkMonitorStatus 获取网络监控状态（详细信息）
func (a *Router) NetworkMonitorStatus(c *fiber.Ctx) error {
	// 检查功能是否启用
	featureStatus := a.NetworkTaskBiz.GetTaskFeatureStatus()
	if !featureStatus["feature_enabled"].(bool) {
		return response.OkWithData(c, map[string]interface{}{
			"feature_enabled": false,
			"message":         "网络监控功能未启用，请在配置文件中启用该功能",
		})
	}

	// 获取详细状态信息
	status, err := a.NetworkTaskBiz.GetTaskStatus()
	if err != nil {
		a.XLog.Errorf("获取网络监控状态失败: %v", err)
		return response.Fail(c, "获取监控状态失败")
	}

	a.XLog.Debugf("网络监控状态查询完成")
	return response.OkWithData(c, status)
}

// NetworkMonitorControl 控制网络监控任务
func (a *Router) NetworkMonitorControl(c *fiber.Ctx) error {
	var req vmodel.NetworkMonitorRequest
	if err := parsex.ParseAll(c, &req); err != nil {
		a.XLog.Errorf("解析网络监控请求参数失败: %v", err)
		return response.Fail(c, "请求参数无效")
	}

	a.XLog.Infof("网络监控控制请求: action=%s", req.Action)

	// 检查功能是否启用（除了查询状态操作）
	if req.Action != "status" {
		featureStatus := a.NetworkTaskBiz.GetTaskFeatureStatus()
		if !featureStatus["feature_enabled"].(bool) {
			a.XLog.Errorf("网络监控功能未启用，无法执行操作: %s", req.Action)
			return response.Fail(c, "网络监控功能未启用，请在配置文件中启用该功能后重启服务")
		}
	}

	var err error
	var message string

	switch req.Action {
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

	case "reset":
		err = a.NetworkTaskBiz.ResetTask()
		if err != nil {
			a.XLog.Errorf("重置网络监控状态失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络监控状态已重置"
		a.XLog.Info(message)

	case "run_now":
		err = a.NetworkTaskBiz.RunTaskNow()
		if err != nil {
			a.XLog.Errorf("立即执行网络检测失败: %v", err)
			return response.Fail(c, err.Error())
		}
		message = "网络检测任务已立即执行"
		a.XLog.Info(message)

	default:
		a.XLog.Errorf("无效的监控操作: %s", req.Action)
		return response.Fail(c, "无效的操作类型，支持的操作: start, stop, restart, reset, run_now")
	}

	return response.OkWithData(c, map[string]interface{}{
		"action":  req.Action,
		"message": message,
	})
}
