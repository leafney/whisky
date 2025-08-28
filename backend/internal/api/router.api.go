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
		a.XLog.Errorf("解析 body 参数操作异常", err)
		return response.Fail(c, "Invalid request body")
	}

	a.XLog.Info(data)

	if status, ok := data[vars.RouterStatus]; ok && status == vars.RouterStsRestart {
		a.XLog.Infof("status %v", status)
		if err := a.RouterSvc.RouterRestart(); err != nil {
			a.XLog.Errorf("RouterStatus error [%v]", err)
			return response.Fail(c, err.Error())
		}
	} else {
		a.XLog.Error("参数错误")
		return response.Fail(c, "参数错误")
	}

	return response.Ok(c)
}

// NetworkMonitorStatus 获取网络监控状态
func (a *Router) NetworkMonitorStatus(c *fiber.Ctx) error {
	status, err := a.NetworkTaskBiz.GetTaskStatus()
	if err != nil {
		a.XLog.Errorf("获取网络监控状态失败: %v", err)
		return response.Fail(c, "获取监控状态失败")
	}
	return response.OkWithData(c, status)
}

// NetworkMonitorControl 控制网络监控
func (a *Router) NetworkMonitorControl(c *fiber.Ctx) error {
	var req vmodel.NetworkMonitorRequest
	if err := parsex.ParseAll(c, &req); err != nil {
		a.XLog.Errorf("解析网络监控请求参数失败: %v", err)
		return response.Fail(c, "请求参数无效")
	}

	a.XLog.Infof("网络监控控制请求: %s", req.Action)

	var err error
	switch req.Action {
	case "start":
		err = a.NetworkTaskBiz.StartTask()
		if err != nil {
			a.XLog.Errorf("启动网络监控失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控已启动")

	case "stop":
		err = a.NetworkTaskBiz.StopTask()
		if err != nil {
			a.XLog.Errorf("停止网络监控失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控已停止")

	case "restart":
		err = a.NetworkTaskBiz.RestartTask()
		if err != nil {
			a.XLog.Errorf("重启网络监控失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控已重启")

	case "reset":
		err = a.NetworkTaskBiz.ResetTask()
		if err != nil {
			a.XLog.Errorf("重置网络监控失败: %v", err)
			return response.Fail(c, err.Error())
		}
		a.XLog.Info("网络监控状态已重置")

	default:
		a.XLog.Errorf("无效的监控操作: %s", req.Action)
		return response.Fail(c, "无效的操作类型")
	}

	return response.Ok(c)
}
