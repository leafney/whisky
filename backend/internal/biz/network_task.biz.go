/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-28
 * @Description: 网络监控任务业务逻辑层
 */

package biz

import (
	"fmt"

	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/pkg/cronx"
	"github.com/leafney/whisky/pkg/xlogx"
)

type NetworkTask struct {
	XLog       *xlogx.XLogSvc
	Config     *config.Config
	CronSvc    *cronx.CronSvc
	MonitorBiz *Monitor
}

const (
	NetworkMonitorTaskID = "network_monitor"
)

// GetTaskStatus 获取网络监控任务状态
func (nt *NetworkTask) GetTaskStatus() (map[string]interface{}, error) {
	// 获取任务运行状态
	isRunning := nt.CronSvc.IsJobExists(NetworkMonitorTaskID)

	// 获取监控详细状态
	status, err := nt.MonitorBiz.GetStatus()
	if err != nil {
		nt.XLog.Errorf("获取网络监控状态失败: %v", err)
		return nil, fmt.Errorf("获取监控状态失败: %v", err)
	}

	stats, err := nt.MonitorBiz.GetStats()
	if err != nil {
		nt.XLog.Errorf("获取网络监控统计失败: %v", err)
		return nil, fmt.Errorf("获取监控统计失败: %v", err)
	}

	result := map[string]interface{}{
		"task_id":         NetworkMonitorTaskID,
		"is_running":      isRunning,
		"feature_enabled": nt.Config.NetworkMonitor.Enable,
		"config_valid":    nt.validateNetworkMonitorConfig() == nil,
		"status":          status,
		"stats":           stats,
	}

	return result, nil
}

// StartTask 启动网络监控任务
func (nt *NetworkTask) StartTask() error {
	// 检查功能是否启用
	if !nt.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，请先在配置文件中设置 Enable=true")
	}

	// 验证配置参数的有效性
	if err := nt.validateNetworkMonitorConfig(); err != nil {
		return fmt.Errorf("网络监控配置无效: %v", err)
	}

	// 如果任务已存在，先停止
	if nt.CronSvc.IsJobExists(NetworkMonitorTaskID) {
		if err := nt.stopTaskInternal(); err != nil {
			nt.XLog.Errorf("停止现有网络监控任务失败: %v", err)
		}
	}

	// 初始化监控状态（从配置文件读取配置）
	if err := nt.MonitorBiz.InitNetworkMonitorFromConfig(); err != nil {
		return fmt.Errorf("初始化网络监控失败: %v", err)
	}

	// 从配置文件构建 cron 表达式（每 N 秒执行一次）
	cronExpr := fmt.Sprintf("*/%d * * * * *", nt.Config.NetworkMonitor.CheckInterval)

	// 注册任务方法
	nt.CronSvc.RegisterTaskMethod("networkMonitor", "网络连通性监控", nt.MonitorBiz.NetworkMonitorJob)

	// 添加定时任务
	if err := nt.CronSvc.AddJobSecs(NetworkMonitorTaskID, cronExpr, nt.MonitorBiz.NetworkMonitorJob); err != nil {
		return fmt.Errorf("添加网络监控定时任务失败: %v", err)
	}

	// 更新监控状态
	if err := nt.MonitorBiz.UpdateStatus(true, "running"); err != nil {
		nt.XLog.Errorf("更新监控状态失败: %v", err)
	}

	nt.XLog.Infof("网络监控任务已启动，检测间隔: %d秒", nt.Config.NetworkMonitor.CheckInterval)
	return nil
}

// StopTask 停止网络监控任务
func (nt *NetworkTask) StopTask() error {
	// 检查功能是否启用
	if !nt.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，无法操作任务")
	}

	return nt.stopTaskInternal()
}

// stopTaskInternal 内部停止网络监控任务的实现
func (nt *NetworkTask) stopTaskInternal() error {
	if !nt.CronSvc.IsJobExists(NetworkMonitorTaskID) {
		return fmt.Errorf("网络监控任务不存在")
	}

	// 移除定时任务
	if err := nt.CronSvc.RemoveJob(NetworkMonitorTaskID); err != nil {
		return fmt.Errorf("移除网络监控任务失败: %v", err)
	}

	// 更新监控状态
	if err := nt.MonitorBiz.UpdateStatus(false, "disabled"); err != nil {
		nt.XLog.Errorf("更新监控状态失败: %v", err)
	}

	nt.XLog.Info("网络监控任务已停止")
	return nil
}

// RestartTask 重启网络监控任务
func (nt *NetworkTask) RestartTask() error {
	// 检查功能是否启用
	if !nt.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，无法操作任务")
	}

	// 先停止
	if err := nt.stopTaskInternal(); err != nil {
		nt.XLog.Errorf("停止网络监控任务失败: %v", err)
	}

	// 再启动
	return nt.StartTask()
}

// RunTaskNow 立即执行一次网络监控
func (nt *NetworkTask) RunTaskNow() error {
	// 检查功能是否启用
	if !nt.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，无法操作任务")
	}

	if !nt.CronSvc.IsJobExists(NetworkMonitorTaskID) {
		return fmt.Errorf("网络监控任务不存在")
	}

	return nt.CronSvc.RunJobNow(NetworkMonitorTaskID)
}

// ResetTask 重置网络监控状态
func (nt *NetworkTask) ResetTask() error {
	if err := nt.MonitorBiz.Reset(); err != nil {
		nt.XLog.Errorf("重置网络监控失败: %v", err)
		return err
	}

	nt.XLog.Info("网络监控状态已重置")
	return nil
}

// IsTaskRunning 检查网络监控任务是否在运行
func (nt *NetworkTask) IsTaskRunning() bool {
	return nt.CronSvc.IsJobExists(NetworkMonitorTaskID)
}

// GetTaskFeatureStatus 获取网络监控功能和任务的完整状态
func (nt *NetworkTask) GetTaskFeatureStatus() map[string]interface{} {
	return map[string]interface{}{
		"feature_enabled": nt.Config.NetworkMonitor.Enable,
		"task_running":    nt.IsTaskRunning(),
		"config_valid":    nt.validateNetworkMonitorConfig() == nil,
	}
}

// validateNetworkMonitorConfig 验证网络监控配置的有效性
func (nt *NetworkTask) validateNetworkMonitorConfig() error {
	config := &nt.Config.NetworkMonitor

	// 验证检测间隔
	if config.CheckInterval <= 0 {
		return fmt.Errorf("检测间隔必须大于0秒，当前值: %d", config.CheckInterval)
	}

	// 验证测试主机列表
	if len(config.TestHosts) == 0 {
		return fmt.Errorf("测试主机列表不能为空")
	}

	// 其他参数使用默认值，无需验证
	// CheckTimeout: 10, FailThreshold: 3, FailHostThreshold: 3, MaxRestarts: 5 等
	// 这些参数在代码中设置为合理的默认值

	return nil
}
