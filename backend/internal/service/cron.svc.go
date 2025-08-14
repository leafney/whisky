package service

import (
	"context"
	"fmt"

	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/pkg/cronx"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Cron struct {
	XLog       *xlogx.XLogSvc
	Config     *config.Config
	CronSvc    *cronx.CronSvc
	MonitorBiz *biz.Monitor
}

const (
	NetworkMonitorTaskID = "network_monitor"
)

func (s *Cron) Start() {
	s.CronSvc.Start()
	s.XLog.Info("定时任务启动成功")
}

func (s *Cron) Stop() {
	s.CronSvc.Stop()
	s.XLog.Info("定时任务停止成功")
}

func (s *Cron) LoadJobs(ctx context.Context) error {
	// 注册任务方法
	s.registerTaskMethods()

	// 这里可以从数据库加载其他定时任务
	// 目前主要是网络监控任务，通过专门的方法管理

	s.CronSvc.AddJobSecs("testJob", "*/5 * * * * *", s.TestJob)

	return nil
}

// GetAvailableMethods 获取所有可用的任务方法
func (s *Cron) GetAvailableMethods() []string {
	methods := make([]string, 0)
	for name := range s.CronSvc.GetRegisteredMethods() {
		methods = append(methods, name)
	}
	return methods
}

// registerTaskMethods 注册所有可用的任务方法
func (s *Cron) registerTaskMethods() {
	s.CronSvc.RegisterTaskMethod("testJob", "测试任务", s.TestJob)
	s.CronSvc.RegisterTaskMethod("networkMonitor", "网络连通性监控", s.MonitorBiz.NetworkMonitorJob)
	// s.CronSvc.RegisterTaskMethod("checkExpired", "检测是否到期", s.CheckExpired)
}

// 测试任务
func (s *Cron) TestJob(ctx context.Context) {
	s.XLog.Info("TestJob 执行了额")
}

// 检测是否到期
func (s *Cron) CheckExpired(ctx context.Context) {
	// s.XLog.Info("[定时任务] -- 检测是否到期 -- 开始执行", zap.String("function", "CheckExpired"))
	// s.CronTaskBiz.CheckExpired(ctx)
}

// ========== 网络监控任务控制方法 ==========

// IsNetworkMonitorFeatureEnabled 检查网络监控功能是否在配置中启用
func (s *Cron) IsNetworkMonitorFeatureEnabled() bool {
	return s.Config.NetworkMonitor.Enable
}

// GetNetworkMonitorFeatureStatus 获取网络监控功能和任务的完整状态
func (s *Cron) GetNetworkMonitorFeatureStatus() map[string]interface{} {
	return map[string]interface{}{
		"feature_enabled": s.Config.NetworkMonitor.Enable,
		"task_running":    s.IsNetworkMonitorRunning(),
		"config_valid":    s.validateNetworkMonitorConfig() == nil,
	}
}

// StartNetworkMonitor 启动网络监控任务（手动控制，需要功能已启用）
func (s *Cron) StartNetworkMonitor() error {
	// 检查功能是否启用（这是前提条件）
	if !s.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，请先在配置文件中设置 Enable=true")
	}

	return s.startNetworkMonitorInternal()
}

// startNetworkMonitorInternal 内部启动网络监控任务的实现
func (s *Cron) startNetworkMonitorInternal() error {
	// 验证配置参数的有效性
	if err := s.validateNetworkMonitorConfig(); err != nil {
		return fmt.Errorf("网络监控配置无效: %v", err)
	}

	// 如果任务已存在，先停止
	if s.CronSvc.IsJobExists(NetworkMonitorTaskID) {
		if err := s.stopNetworkMonitorInternal(); err != nil {
			s.XLog.Errorf("停止现有网络监控任务失败: %v", err)
		}
	}

	// 初始化监控状态（从配置文件读取配置）
	if err := s.MonitorBiz.InitNetworkMonitorFromConfig(); err != nil {
		return fmt.Errorf("初始化网络监控失败: %v", err)
	}

	// 从配置文件构建 cron 表达式（每 N 秒执行一次）
	cronExpr := fmt.Sprintf("*/%d * * * * *", s.Config.NetworkMonitor.CheckInterval)

	// 添加定时任务
	if err := s.CronSvc.AddJobSecs(NetworkMonitorTaskID, cronExpr, s.MonitorBiz.NetworkMonitorJob); err != nil {
		return fmt.Errorf("添加网络监控定时任务失败: %v", err)
	}

	// 更新监控状态
	if err := s.MonitorBiz.UpdateStatus(true, "running"); err != nil {
		s.XLog.Errorf("更新监控状态失败: %v", err)
	}

	s.XLog.Infof("网络监控任务已启动，检测间隔: %d秒", s.Config.NetworkMonitor.CheckInterval)
	return nil
}

// validateNetworkMonitorConfig 验证网络监控配置的有效性
func (s *Cron) validateNetworkMonitorConfig() error {
	config := &s.Config.NetworkMonitor

	if config.CheckInterval <= 0 {
		return fmt.Errorf("检测间隔必须大于0秒，当前值: %d", config.CheckInterval)
	}

	if config.CheckTimeout <= 0 {
		return fmt.Errorf("检测超时必须大于0秒，当前值: %d", config.CheckTimeout)
	}

	if len(config.TestHosts) == 0 {
		return fmt.Errorf("测试主机列表不能为空")
	}

	if config.FailThreshold <= 0 {
		return fmt.Errorf("失败阈值必须大于0，当前值: %d", config.FailThreshold)
	}

	if config.FailHostThreshold <= 0 || config.FailHostThreshold > len(config.TestHosts) {
		return fmt.Errorf("失败主机阈值必须在1到%d之间，当前值: %d", len(config.TestHosts), config.FailHostThreshold)
	}

	if config.MaxRestarts < 0 {
		return fmt.Errorf("最大重启次数不能为负数，当前值: %d", config.MaxRestarts)
	}

	return nil
}

// StopNetworkMonitor 停止网络监控任务（手动控制，需要功能已启用）
func (s *Cron) StopNetworkMonitor() error {
	// 检查功能是否启用（这是前提条件）
	if !s.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，无法操作任务")
	}

	return s.stopNetworkMonitorInternal()
}

// stopNetworkMonitorInternal 内部停止网络监控任务的实现
func (s *Cron) stopNetworkMonitorInternal() error {
	if !s.CronSvc.IsJobExists(NetworkMonitorTaskID) {
		return fmt.Errorf("网络监控任务不存在")
	}

	// 移除定时任务
	if err := s.CronSvc.RemoveJob(NetworkMonitorTaskID); err != nil {
		return fmt.Errorf("移除网络监控任务失败: %v", err)
	}

	// 更新监控状态
	if err := s.MonitorBiz.UpdateStatus(false, "disabled"); err != nil {
		s.XLog.Errorf("更新监控状态失败: %v", err)
	}

	s.XLog.Info("网络监控任务已停止")
	return nil
}

// RestartNetworkMonitor 重启网络监控任务
func (s *Cron) RestartNetworkMonitor() error {
	// 检查功能是否启用（这是前提条件）
	if !s.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，无法操作任务")
	}

	// 先停止（使用内部方法，避免重复检查）
	if err := s.stopNetworkMonitorInternal(); err != nil {
		s.XLog.Errorf("停止网络监控任务失败: %v", err)
	}

	// 再启动（使用内部方法，避免重复检查）
	return s.startNetworkMonitorInternal()
}

// IsNetworkMonitorRunning 检查网络监控任务是否在运行
func (s *Cron) IsNetworkMonitorRunning() bool {
	return s.CronSvc.IsJobExists(NetworkMonitorTaskID)
}

// RunNetworkMonitorNow 立即执行一次网络监控
func (s *Cron) RunNetworkMonitorNow() error {
	// 检查功能是否启用（这是前提条件）
	if !s.Config.NetworkMonitor.Enable {
		return fmt.Errorf("网络监控功能未启用，无法操作任务")
	}

	if !s.CronSvc.IsJobExists(NetworkMonitorTaskID) {
		return fmt.Errorf("网络监控任务不存在")
	}

	return s.CronSvc.RunJobNow(NetworkMonitorTaskID)
}

// AutoStartNetworkMonitor 自动启动网络监控（仅在配置启用时启动）
func (s *Cron) AutoStartNetworkMonitor() error {
	// 检查功能是否启用
	if !s.Config.NetworkMonitor.Enable {
		s.XLog.Info("网络监控功能未启用，跳过自动启动")
		return nil
	}

	// 功能已启用，执行自动启动
	if err := s.startNetworkMonitorInternal(); err != nil {
		s.XLog.Errorf("自动启动网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控已自动启动")
	return nil
}
