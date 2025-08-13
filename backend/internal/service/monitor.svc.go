/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 网络监控服务层（基于定时任务）
 */

package service

import (
	"context"
	"fmt"

	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/xlogx"
)

type NetworkMonitor struct {
	XLog       *xlogx.XLogSvc
	Config     *config.Config
	MonitorBiz *biz.Monitor
}

// Initialize 初始化网络监控服务
func (s *NetworkMonitor) Initialize() error {
	// 初始化定时任务调度器
	if err := s.MonitorBiz.InitScheduler(); err != nil {
		s.XLog.Errorf("初始化网络监控调度器失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控服务初始化完成")
	return nil
}

// AutoStart 自动启动监控（如果配置启用）
func (s *NetworkMonitor) AutoStart(ctx context.Context) error {
	// 检查配置是否启用了网络监控
	if !s.Config.NetworkMonitor.Enable {
		s.XLog.Info("网络监控功能未启用")
		return nil
	}

	// 将配置转换为vmodel格式
	config := &vmodel.NetworkMonitorConfig{
		Enable:            s.Config.NetworkMonitor.Enable,
		CheckInterval:     s.Config.NetworkMonitor.CheckInterval,
		FailCheckInterval: s.Config.NetworkMonitor.FailCheckInterval,
		CheckTimeout:      s.Config.NetworkMonitor.CheckTimeout,
		TestHosts:         s.Config.NetworkMonitor.TestHosts,
		FailThreshold:     s.Config.NetworkMonitor.FailThreshold,
		FailHostThreshold: s.Config.NetworkMonitor.FailHostThreshold,
		MaxRestarts:       s.Config.NetworkMonitor.MaxRestarts,
		RestartWindow:     s.Config.NetworkMonitor.RestartWindow,
		CooldownPeriod:    s.Config.NetworkMonitor.CooldownPeriod,
	}

	// 启动监控
	if err := s.StartMonitor(config); err != nil {
		s.XLog.Errorf("自动启动网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控已自动启动")
	return nil
}

// GetMonitorStatus 获取网络监控状态
func (s *NetworkMonitor) GetMonitorStatus() (*vmodel.NetworkMonitorStatus, error) {
	status, err := s.MonitorBiz.GetStatus()
	if err != nil {
		s.XLog.Errorf("获取监控状态失败: %v", err)
		return nil, err
	}

	return status, nil
}

// GetMonitorStats 获取网络监控统计信息
func (s *NetworkMonitor) GetMonitorStats() (*vmodel.NetworkMonitorStats, error) {
	stats, err := s.MonitorBiz.GetStats()
	if err != nil {
		s.XLog.Errorf("获取监控统计失败: %v", err)
		return nil, err
	}

	return stats, nil
}

// GetDetailedStatus 获取详细的监控状态（包括统计信息）
func (s *NetworkMonitor) GetDetailedStatus() (map[string]interface{}, error) {
	status, err := s.GetMonitorStatus()
	if err != nil {
		return nil, err
	}

	stats, err := s.GetMonitorStats()
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"status": status,
		"stats":  stats,
	}

	return result, nil
}

// StartMonitor 启动网络监控定时任务
func (s *NetworkMonitor) StartMonitor(config *vmodel.NetworkMonitorConfig) error {
	// 如果没有提供配置，使用默认配置
	if config == nil {
		config = s.getDefaultConfig()
	}

	// 验证配置
	if err := s.validateConfig(config); err != nil {
		s.XLog.Errorf("配置验证失败: %v", err)
		return fmt.Errorf("配置无效: %v", err)
	}

	// 启动监控
	if err := s.MonitorBiz.Start(config); err != nil {
		s.XLog.Errorf("启动网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控定时任务已启动")
	return nil
}

// StopMonitor 停止网络监控定时任务
func (s *NetworkMonitor) StopMonitor() error {
	if err := s.MonitorBiz.Stop(); err != nil {
		s.XLog.Errorf("停止网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控定时任务已停止")
	return nil
}

// RestartMonitor 重启网络监控定时任务
func (s *NetworkMonitor) RestartMonitor(config *vmodel.NetworkMonitorConfig) error {
	// 如果没有提供配置，使用当前配置
	if config == nil {
		currentStatus, err := s.GetMonitorStatus()
		if err != nil {
			return fmt.Errorf("获取当前配置失败: %v", err)
		}
		config = &currentStatus.Config
	}

	// 验证配置
	if err := s.validateConfig(config); err != nil {
		return fmt.Errorf("配置无效: %v", err)
	}

	// 重启监控
	if err := s.MonitorBiz.Restart(config); err != nil {
		s.XLog.Errorf("重启网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控定时任务已重启")
	return nil
}

// ResetMonitor 重置网络监控状态
func (s *NetworkMonitor) ResetMonitor() error {
	if err := s.MonitorBiz.Reset(); err != nil {
		s.XLog.Errorf("重置网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控状态已重置")
	return nil
}

// UpdateMonitorConfig 更新网络监控配置
func (s *NetworkMonitor) UpdateMonitorConfig(config *vmodel.NetworkMonitorConfig) error {
	// 验证配置
	if err := s.validateConfig(config); err != nil {
		s.XLog.Errorf("配置验证失败: %v", err)
		return fmt.Errorf("配置无效: %v", err)
	}

	// 更新配置
	if err := s.MonitorBiz.UpdateConfig(config); err != nil {
		s.XLog.Errorf("更新监控配置失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控配置已更新")
	return nil
}

// Shutdown 优雅关闭网络监控服务
func (s *NetworkMonitor) Shutdown() error {
	// 先停止监控任务
	if err := s.StopMonitor(); err != nil {
		s.XLog.Errorf("停止监控任务失败: %v", err)
	}

	// 关闭调度器
	if err := s.MonitorBiz.Shutdown(); err != nil {
		s.XLog.Errorf("关闭调度器失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控服务已优雅关闭")
	return nil
}

// 获取默认配置
func (s *NetworkMonitor) getDefaultConfig() *vmodel.NetworkMonitorConfig {
	return &vmodel.NetworkMonitorConfig{
		Enable:            true,
		CheckInterval:     300, // 5分钟
		FailCheckInterval: 60,  // 1分钟
		CheckTimeout:      10,  // 10秒
		TestHosts: []string{
			"8.8.8.8",
			"114.114.114.114",
			"1.1.1.1",
			"223.5.5.5",
			"www.baidu.com",
		},
		FailThreshold:     3,  // 连续3次失败
		FailHostThreshold: 3,  // 3个主机失败
		MaxRestarts:       5,  // 最多5次重启
		RestartWindow:     24, // 24小时窗口
		CooldownPeriod:    30, // 30分钟冷却
	}
}

// 验证配置
func (s *NetworkMonitor) validateConfig(config *vmodel.NetworkMonitorConfig) error {
	if config == nil {
		return fmt.Errorf("配置不能为空")
	}

	if config.CheckInterval < 30 {
		return fmt.Errorf("检测间隔不能少于30秒")
	}

	if config.FailCheckInterval < 10 {
		return fmt.Errorf("失败检测间隔不能少于10秒")
	}

	if config.CheckTimeout < 1 || config.CheckTimeout > 60 {
		return fmt.Errorf("检测超时时间必须在1-60秒之间")
	}

	if len(config.TestHosts) == 0 {
		return fmt.Errorf("测试主机列表不能为空")
	}

	if len(config.TestHosts) > 20 {
		return fmt.Errorf("测试主机数量不能超过20个")
	}

	if config.FailThreshold < 1 || config.FailThreshold > 10 {
		return fmt.Errorf("失败阈值必须在1-10次之间")
	}

	if config.FailHostThreshold < 1 || config.FailHostThreshold > len(config.TestHosts) {
		return fmt.Errorf("失败主机阈值必须在1-%d之间", len(config.TestHosts))
	}

	if config.MaxRestarts < 1 || config.MaxRestarts > 20 {
		return fmt.Errorf("最大重启次数必须在1-20次之间")
	}

	if config.RestartWindow < 1 || config.RestartWindow > 168 {
		return fmt.Errorf("重启窗口期必须在1-168小时之间")
	}

	if config.CooldownPeriod < 5 || config.CooldownPeriod > 1440 {
		return fmt.Errorf("冷却期必须在5-1440分钟之间")
	}

	return nil
}
