/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 网络监控业务逻辑层（纯检测逻辑）
 */

package biz

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/internal/dao"
	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/cmds"
	"github.com/leafney/whisky/pkg/utils"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Monitor struct {
	XLog       *xlogx.XLogSvc
	Config     *config.Config
	MonitorDao *dao.Monitor
}

// 全局监控状态（简化状态管理）
var (
	monitorStatus *vmodel.NetworkMonitorStatus
	monitorStats  *vmodel.NetworkMonitorStats
	stateMutex    sync.RWMutex
)

// GetStatus 获取监控状态
func (b *Monitor) GetStatus() (*vmodel.NetworkMonitorStatus, error) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()

	if monitorStatus == nil {
		// 从数据库加载状态
		status, err := b.MonitorDao.GetMonitorStatus()
		if err != nil {
			return nil, err
		}
		monitorStatus = status
	}

	return monitorStatus, nil
}

// GetStats 获取监控统计
func (b *Monitor) GetStats() (*vmodel.NetworkMonitorStats, error) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()

	if monitorStats == nil {
		stats, err := b.MonitorDao.GetMonitorStats()
		if err != nil {
			return nil, err
		}
		monitorStats = stats
	}

	return monitorStats, nil
}

// InitializeStatus 初始化监控状态
func (b *Monitor) InitializeStatus(config *vmodel.NetworkMonitorConfig) error {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	// 尝试从数据库加载现有状态
	status, err := b.MonitorDao.GetMonitorStatus()
	if err != nil {
		return err
	}

	// 如果没有现有状态，创建新状态
	if status == nil {
		now := time.Now()
		status = &vmodel.NetworkMonitorStatus{
			Enabled:          true,
			CurrentStatus:    "running",
			ConsecutiveFails: 0,
			TotalRestarts:    0,
			RestartsInWindow: 0,
			WindowStartTime:  &now,
			Config:           *config,
		}
	} else {
		// 更新配置
		status.Config = *config
		status.Enabled = true
		status.CurrentStatus = "running"
	}

	monitorStatus = status

	// 初始化统计
	stats, err := b.MonitorDao.GetMonitorStats()
	if err != nil {
		return err
	}
	if stats == nil {
		stats = &vmodel.NetworkMonitorStats{}
	}
	monitorStats = stats

	return nil
}

// UpdateStatus 更新监控状态
func (b *Monitor) UpdateStatus(enabled bool, status string) error {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	if monitorStatus != nil {
		monitorStatus.Enabled = enabled
		monitorStatus.CurrentStatus = status
		return b.MonitorDao.SaveMonitorStatus(monitorStatus)
	}
	return nil
}

// PerformNetworkCheck 执行网络检测（供定时任务调用）
func (b *Monitor) PerformNetworkCheck(ctx context.Context) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	if monitorStatus == nil {
		b.XLog.Error("监控状态未初始化")
		return
	}

	// 检查是否在冷却期
	if b.isInCooldown() {
		b.XLog.Debug("处于冷却期，跳过检测")
		return
	}

	b.XLog.Debug("开始执行网络连通性检测")
	now := time.Now()
	monitorStatus.LastCheckTime = &now

	// 并发检测所有主机
	results := b.checkAllHosts(ctx)
	monitorStatus.LastCheckResults = results

	// 分析检测结果
	failedHosts := 0
	for _, result := range results {
		if !result.Success {
			failedHosts++
		}
	}

	// 更新统计
	monitorStats.TotalChecks++

	// 判断本次检测是否失败
	checkFailed := failedHosts >= monitorStatus.Config.FailHostThreshold

	if checkFailed {
		monitorStatus.ConsecutiveFails++
		monitorStats.FailedChecks++
		b.XLog.Errorf("网络检测失败，连续失败次数: %d/%d，失败主机: %d/%d",
			monitorStatus.ConsecutiveFails, monitorStatus.Config.FailThreshold,
			failedHosts, len(monitorStatus.Config.TestHosts))

		// 检查是否需要重启
		if b.shouldRestart() {
			b.executeRestart()
		}
	} else {
		// 检测成功，重置失败计数
		if monitorStatus.ConsecutiveFails > 0 {
			b.XLog.Info("网络连通性恢复正常")
			monitorStatus.ConsecutiveFails = 0
		}
		monitorStats.SuccessfulChecks++
	}

	// 更新成功率
	monitorStats.SuccessRate = float64(monitorStats.SuccessfulChecks) / float64(monitorStats.TotalChecks) * 100

	// 计算下次检测时间
	var nextInterval time.Duration
	if checkFailed && monitorStatus.ConsecutiveFails < monitorStatus.Config.FailThreshold {
		nextInterval = time.Duration(monitorStatus.Config.FailCheckInterval) * time.Second
	} else {
		nextInterval = time.Duration(monitorStatus.Config.CheckInterval) * time.Second
	}
	nextCheck := now.Add(nextInterval)
	monitorStatus.NextCheckTime = &nextCheck

	// 保存状态和统计
	b.MonitorDao.SaveMonitorStatus(monitorStatus)
	b.MonitorDao.SaveMonitorStats(monitorStats)
}

// Reset 重置监控状态
func (b *Monitor) Reset() error {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	if monitorStatus != nil {
		monitorStatus.ConsecutiveFails = 0
		monitorStatus.TotalRestarts = 0
		monitorStatus.RestartsInWindow = 0
		monitorStatus.LastRestartTime = nil
		now := time.Now()
		monitorStatus.WindowStartTime = &now

		// 保存状态
		if err := b.MonitorDao.SaveMonitorStatus(monitorStatus); err != nil {
			return fmt.Errorf("保存重置状态失败: %v", err)
		}
	}

	// 重置统计
	if monitorStats != nil {
		monitorStats.TotalChecks = 0
		monitorStats.SuccessfulChecks = 0
		monitorStats.FailedChecks = 0
		monitorStats.SuccessRate = 0.0

		if err := b.MonitorDao.SaveMonitorStats(monitorStats); err != nil {
			return fmt.Errorf("保存重置统计失败: %v", err)
		}
	}

	b.XLog.Info("网络监控状态已重置")
	return nil
}

// 检测所有主机
func (b *Monitor) checkAllHosts(ctx context.Context) []vmodel.HostCheckResult {
	hosts := monitorStatus.Config.TestHosts
	results := make([]vmodel.HostCheckResult, len(hosts))
	var wg sync.WaitGroup

	for i, host := range hosts {
		wg.Add(1)
		go func(index int, hostname string) {
			defer wg.Done()
			results[index] = b.pingHost(ctx, hostname)
		}(i, host)
	}

	wg.Wait()
	return results
}

// pingHost 检测单个主机连通性
func (b *Monitor) pingHost(ctx context.Context, host string) vmodel.HostCheckResult {
	start := time.Now()
	result := vmodel.HostCheckResult{
		Host:      host,
		CheckTime: start,
	}

	// 设置超时
	timeout := time.Duration(monitorStatus.Config.CheckTimeout) * time.Second

	// 根据主机类型选择检测策略
	if b.isWebHost(host) {
		// 对于网站域名，优先使用HTTP检测
		if b.checkHTTP(host, timeout, &result) {
			result.Latency = time.Since(start)
			return result
		}
	} else if b.isIPAddress(host) {
		// 对于IP地址，根据类型选择检测方式
		if b.isDNSServer(host) {
			// DNS服务器使用UDP 53端口检测
			if b.checkDNS(host, timeout, &result) {
				result.Latency = time.Since(start)
				return result
			}
		}
	}

	// 回退到ping检测
	if b.checkPing(ctx, host, timeout, &result) {
		result.Latency = time.Since(start)
		return result
	}

	// 所有检测方式都失败
	result.Success = false
	if result.Error == "" {
		result.Error = "所有连通性检测方式都失败"
	}
	result.Latency = time.Since(start)
	return result
}

// isWebHost 判断是否为网站域名
func (b *Monitor) isWebHost(host string) bool {
	return !b.isIPAddress(host) && (host == "www.baidu.com" || host == "www.google.com" ||
		host == "www.bing.com" || host == "www.qq.com")
}

// isIPAddress 判断是否为IP地址
func (b *Monitor) isIPAddress(host string) bool {
	return net.ParseIP(host) != nil
}

// isDNSServer 判断是否为知名DNS服务器
func (b *Monitor) isDNSServer(host string) bool {
	dnsServers := []string{"8.8.8.8", "8.8.4.4", "114.114.114.114", "223.5.5.5", "1.1.1.1", "208.67.222.222"}
	for _, dns := range dnsServers {
		if host == dns {
			return true
		}
	}
	return false
}

// checkHTTP HTTP连通性检测
func (b *Monitor) checkHTTP(host string, timeout time.Duration, result *vmodel.HostCheckResult) bool {
	// 先尝试HTTPS 443端口
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "443"), timeout); err == nil {
		conn.Close()
		result.Success = true
		return true
	}

	// 再尝试HTTP 80端口
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), timeout); err == nil {
		conn.Close()
		result.Success = true
		return true
	}

	return false
}

// checkDNS DNS服务器连通性检测
func (b *Monitor) checkDNS(host string, timeout time.Duration, result *vmodel.HostCheckResult) bool {
	// 尝试UDP 53端口（DNS标准端口）
	if conn, err := net.DialTimeout("udp", net.JoinHostPort(host, "53"), timeout); err == nil {
		conn.Close()
		result.Success = true
		return true
	}

	// 也可以尝试TCP 53端口
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "53"), timeout); err == nil {
		conn.Close()
		result.Success = true
		return true
	}

	return false
}

// checkPing Ping连通性检测
func (b *Monitor) checkPing(ctx context.Context, host string, timeout time.Duration, result *vmodel.HostCheckResult) bool {
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pingCmd := fmt.Sprintf("ping -c 1 -W %d %s", monitorStatus.Config.CheckTimeout, host)
	if _, err := utils.RunBashCtx(pingCtx, pingCmd); err != nil {
		result.Error = fmt.Sprintf("ping检测失败: %v", err)
		return false
	}

	result.Success = true
	return true
}

// 判断是否应该重启
func (b *Monitor) shouldRestart() bool {
	// 检查连续失败次数
	if monitorStatus.ConsecutiveFails < monitorStatus.Config.FailThreshold {
		return false
	}

	// 检查时间窗口内的重启次数
	if b.isRestartLimitReached() {
		b.XLog.Error("已达到最大重启次数限制，自动禁用监控")
		monitorStatus.Enabled = false
		monitorStatus.CurrentStatus = "disabled"
		return false
	}

	return true
}

// 检查是否达到重启限制
func (b *Monitor) isRestartLimitReached() bool {
	now := time.Now()
	windowHours := time.Duration(monitorStatus.Config.RestartWindow) * time.Hour

	// 如果窗口开始时间为空或超过窗口期，重置窗口
	if monitorStatus.WindowStartTime == nil || now.Sub(*monitorStatus.WindowStartTime) > windowHours {
		monitorStatus.WindowStartTime = &now
		monitorStatus.RestartsInWindow = 0
		return false
	}

	return monitorStatus.RestartsInWindow >= monitorStatus.Config.MaxRestarts
}

// 检查是否在冷却期
func (b *Monitor) isInCooldown() bool {
	if monitorStatus.LastRestartTime == nil {
		return false
	}

	cooldownDuration := time.Duration(monitorStatus.Config.CooldownPeriod) * time.Minute
	return time.Since(*monitorStatus.LastRestartTime) < cooldownDuration
}

// 执行重启操作
func (b *Monitor) executeRestart() {
	b.XLog.Error("网络连通性持续异常，执行路由器重启")

	// 更新重启统计
	now := time.Now()
	monitorStatus.LastRestartTime = &now
	monitorStatus.TotalRestarts++
	monitorStatus.RestartsInWindow++
	monitorStatus.ConsecutiveFails = 0
	monitorStatus.CurrentStatus = "cooldown"

	// 记录重启日志
	reason := fmt.Sprintf("连续%d次网络检测失败", monitorStatus.Config.FailThreshold)
	b.MonitorDao.LogRestart(reason)

	// 保存状态
	b.MonitorDao.SaveMonitorStatus(monitorStatus)

	// 执行重启命令（异步）
	go func() {
		// 延迟执行，确保状态已保存
		time.Sleep(2 * time.Second)
		if _, err := utils.RunBash(cmds.ScriptReboot); err != nil {
			b.XLog.Errorf("执行重启命令失败: %v", err)
		} else {
			b.XLog.Info("路由器重启命令已执行")
		}
	}()
}

// ========== 定时任务相关方法 ==========

// NetworkMonitorJob 网络监控定时任务执行函数
func (b *Monitor) NetworkMonitorJob(ctx context.Context) {
	b.XLog.Debug("执行网络监控定时任务")
	// 调用网络检测逻辑
	b.PerformNetworkCheck(ctx)
}

// InitNetworkMonitor 初始化网络监控（从配置文件读取配置）
func (b *Monitor) InitNetworkMonitorFromConfig() error {
	// 从配置文件获取网络监控配置
	config := &vmodel.NetworkMonitorConfig{
		Enable:            b.Config.NetworkMonitor.Enable,
		CheckInterval:     b.Config.NetworkMonitor.CheckInterval,
		FailCheckInterval: b.Config.NetworkMonitor.FailCheckInterval,
		CheckTimeout:      b.Config.NetworkMonitor.CheckTimeout,
		TestHosts:         b.Config.NetworkMonitor.TestHosts,
		FailThreshold:     b.Config.NetworkMonitor.FailThreshold,
		FailHostThreshold: b.Config.NetworkMonitor.FailHostThreshold,
		MaxRestarts:       b.Config.NetworkMonitor.MaxRestarts,
		RestartWindow:     b.Config.NetworkMonitor.RestartWindow,
		CooldownPeriod:    b.Config.NetworkMonitor.CooldownPeriod,
	}

	// 初始化监控状态
	if err := b.InitializeStatus(config); err != nil {
		b.XLog.Errorf("初始化网络监控状态失败: %v", err)
		return err
	}

	b.XLog.Info("网络监控初始化完成")
	return nil
}
