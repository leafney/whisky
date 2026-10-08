/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 网络监控服务层（简化版，主要用于状态查询）
 */

package service

import (
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/xlogx"
)

type NetworkMonitor struct {
	XLog       *xlogx.XLogSvc
	MonitorBiz *biz.Monitor
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

// ResetMonitor 重置网络监控状态
func (s *NetworkMonitor) ResetMonitor() error {
	if err := s.MonitorBiz.Reset(); err != nil {
		s.XLog.Errorf("重置网络监控失败: %v", err)
		return err
	}

	s.XLog.Info("网络监控状态已重置")
	return nil
}
