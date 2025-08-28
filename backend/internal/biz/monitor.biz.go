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
		b.XLog.Error("监控状态未初始化，无法执行网络检测")
		return
	}

	// 检查是否在冷却期
	if b.isInCooldown() {
		remainingTime := time.Duration(monitorStatus.Config.CooldownPeriod)*time.Minute - time.Since(*monitorStatus.LastRestartTime)
		b.XLog.Debugf("处于冷却期，跳过检测。剩余冷却时间: %v", remainingTime.Round(time.Second))
		return
	}

	b.XLog.Infof("开始执行网络连通性检测 - 连续失败次数: %d/%d", 
		monitorStatus.ConsecutiveFails, monitorStatus.Config.FailThreshold)
	now := time.Now()
	monitorStatus.LastCheckTime = &now

	// 并发检测所有主机
	b.XLog.Debugf("开始检测 %d 个主机: %v", len(monitorStatus.Config.TestHosts), monitorStatus.Config.TestHosts)
	results := b.checkAllHosts(ctx)
	monitorStatus.LastCheckResults = results

	// 分析检测结果
	failedHosts := 0
	successHosts := []string{}
	failedHostList := []string{}
	for _, result := range results {
		if !result.Success {
			failedHosts++
			failedHostList = append(failedHostList, fmt.Sprintf("%s(%s)", result.Host, result.Error))
		} else {
			successHosts = append(successHosts, fmt.Sprintf("%s(%.0fms)", result.Host, float64(result.Latency.Nanoseconds())/1e6))
		}
	}

	b.XLog.Infof("检测结果汇总 - 成功: %d/%d, 失败: %d/%d", 
		len(results)-failedHosts, len(results), failedHosts, len(results))
	if len(successHosts) > 0 {
		b.XLog.Debugf("成功主机: %v", successHosts)
	}
	if len(failedHostList) > 0 {
		b.XLog.Debugf("失败主机: %v", failedHostList)
	}

	// 更新统计
	monitorStats.TotalChecks++

	// 判断本次检测是否失败
	checkFailed := failedHosts >= monitorStatus.Config.FailHostThreshold

	if checkFailed {
		monitorStatus.ConsecutiveFails++
		monitorStats.FailedChecks++
		b.XLog.Errorf("网络检测失败 - 连续失败: %d/%d, 本次失败主机: %d/%d (阈值: %d)",
			monitorStatus.ConsecutiveFails, monitorStatus.Config.FailThreshold,
			failedHosts, len(monitorStatus.Config.TestHosts), monitorStatus.Config.FailHostThreshold)

		// 检查是否需要重启
		if b.shouldRestart() {
			b.XLog.Errorf("网络连通性持续异常，准备执行路由器重启 - 总重启次数: %d, 窗口内重启: %d/%d",
				monitorStatus.TotalRestarts, monitorStatus.RestartsInWindow, monitorStatus.Config.MaxRestarts)
			b.executeRestart()
		} else {
			b.XLog.Infof("网络检测失败，但未达到重启条件 - 需要连续失败 %d 次", 
				monitorStatus.Config.FailThreshold-monitorStatus.ConsecutiveFails)
		}
	} else {
		// 检测成功，重置失败计数
		if monitorStatus.ConsecutiveFails > 0 {
			b.XLog.Infof("网络连通性恢复正常 - 重置连续失败计数 %d -> 0", monitorStatus.ConsecutiveFails)
			monitorStatus.ConsecutiveFails = 0
		} else {
			b.XLog.Debugf("网络连通性正常 - 平均延迟: %.0fms", b.calculateAverageLatency(results))
		}
		monitorStats.SuccessfulChecks++
	}

	// 更新成功率
	monitorStats.SuccessRate = float64(monitorStats.SuccessfulChecks) / float64(monitorStats.TotalChecks) * 100

	// 计算下次检测时间
	var nextInterval time.Duration
	var intervalType string
	if checkFailed && monitorStatus.ConsecutiveFails < monitorStatus.Config.FailThreshold {
		nextInterval = time.Duration(monitorStatus.Config.FailCheckInterval) * time.Second
		intervalType = "失败检测间隔"
	} else {
		nextInterval = time.Duration(monitorStatus.Config.CheckInterval) * time.Second
		intervalType = "正常检测间隔"
	}
	nextCheck := now.Add(nextInterval)
	monitorStatus.NextCheckTime = &nextCheck

	b.XLog.Debugf("本次检测完成 - 成功率: %.1f%% (%d/%d), 下次检测: %s (%s)",
		monitorStats.SuccessRate, monitorStats.SuccessfulChecks, monitorStats.TotalChecks,
		nextCheck.Format("15:04:05"), intervalType)

	// 保存状态和统计
	if err := b.MonitorDao.SaveMonitorStatus(monitorStatus); err != nil {
		b.XLog.Errorf("保存监控状态失败: %v", err)
	}
	if err := b.MonitorDao.SaveMonitorStats(monitorStats); err != nil {
		b.XLog.Errorf("保存监控统计失败: %v", err)
	}
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

	checkStart := time.Now()
	b.XLog.Debugf("开始并发检测 %d 个主机，超时时间: %ds", len(hosts), monitorStatus.Config.CheckTimeout)

	for i, host := range hosts {
		wg.Add(1)
		go func(index int, hostname string) {
			defer wg.Done()
			hostStart := time.Now()
			results[index] = b.pingHost(ctx, hostname)
			b.XLog.Debugf("主机 %s 检测完成: 成功=%v, 耗时=%.0fms", 
				hostname, results[index].Success, float64(time.Since(hostStart).Nanoseconds())/1e6)
		}(i, host)
	}

	wg.Wait()
	totalTime := time.Since(checkStart)
	b.XLog.Debugf("所有主机检测完成，总耗时: %.0fms", float64(totalTime.Nanoseconds())/1e6)
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
	b.XLog.Debugf("开始检测主机 %s (超时: %v)", host, timeout)

	// 根据主机类型选择检测策略
	if b.isWebHost(host) {
		// 对于网站域名，优先使用HTTP检测
		b.XLog.Debugf("检测网站域名 %s - 尝试HTTP连接", host)
		if b.checkHTTP(host, timeout, &result) {
			result.Latency = time.Since(start)
			b.XLog.Debugf("主机 %s HTTP检测成功，延迟: %.0fms", host, float64(result.Latency.Nanoseconds())/1e6)
			return result
		}
		b.XLog.Debugf("主机 %s HTTP检测失败，尝试ping检测", host)
	} else if b.isIPAddress(host) {
		// 对于IP地址，根据类型选择检测方式
		if b.isDNSServer(host) {
			// DNS服务器使用UDP 53端口检测
			b.XLog.Debugf("检测DNS服务器 %s - 尝试UDP 53端口", host)
			if b.checkDNS(host, timeout, &result) {
				result.Latency = time.Since(start)
				b.XLog.Debugf("DNS服务器 %s 检测成功，延迟: %.0fms", host, float64(result.Latency.Nanoseconds())/1e6)
				return result
			}
			b.XLog.Debugf("DNS服务器 %s UDP检测失败，尝试ping检测", host)
		} else {
			b.XLog.Debugf("检测IP地址 %s - 直接使用ping", host)
		}
	}

	// 回退到ping检测
	if b.checkPing(ctx, host, timeout, &result) {
		result.Latency = time.Since(start)
		b.XLog.Debugf("主机 %s ping检测成功，延迟: %.0fms", host, float64(result.Latency.Nanoseconds())/1e6)
		return result
	}

	// 所有检测方式都失败
	result.Success = false
	if result.Error == "" {
		result.Error = "所有连通性检测方式都失败"
	}
	result.Latency = time.Since(start)
	b.XLog.Errorf("主机 %s 所有检测方式都失败，总耗时: %.0fms, 错误: %s", 
		host, float64(result.Latency.Nanoseconds())/1e6, result.Error)
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
	b.XLog.Debugf("尝试HTTPS连接 %s:443", host)
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "443"), timeout); err == nil {
		conn.Close()
		result.Success = true
		b.XLog.Debugf("HTTPS连接成功: %s:443", host)
		return true
	} else {
		b.XLog.Debugf("HTTPS连接失败 %s:443 - %v", host, err)
	}

	// 再尝试HTTP 80端口
	b.XLog.Debugf("尝试HTTP连接 %s:80", host)
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), timeout); err == nil {
		conn.Close()
		result.Success = true
		b.XLog.Debugf("HTTP连接成功: %s:80", host)
		return true
	} else {
		b.XLog.Debugf("HTTP连接失败 %s:80 - %v", host, err)
		result.Error = fmt.Sprintf("HTTP/HTTPS连接失败: %v", err)
	}

	return false
}

// checkDNS DNS服务器连通性检测
func (b *Monitor) checkDNS(host string, timeout time.Duration, result *vmodel.HostCheckResult) bool {
	// 尝试UDP 53端口（DNS标准端口）
	b.XLog.Debugf("尝试UDP DNS连接 %s:53", host)
	if conn, err := net.DialTimeout("udp", net.JoinHostPort(host, "53"), timeout); err == nil {
		conn.Close()
		result.Success = true
		b.XLog.Debugf("UDP DNS连接成功: %s:53", host)
		return true
	} else {
		b.XLog.Debugf("UDP DNS连接失败 %s:53 - %v", host, err)
	}

	// 也可以尝试TCP 53端口
	b.XLog.Debugf("尝试TCP DNS连接 %s:53", host)
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "53"), timeout); err == nil {
		conn.Close()
		result.Success = true
		b.XLog.Debugf("TCP DNS连接成功: %s:53", host)
		return true
	} else {
		b.XLog.Debugf("TCP DNS连接失败 %s:53 - %v", host, err)
		result.Error = fmt.Sprintf("DNS连接失败(UDP+TCP): %v", err)
	}

	return false
}

// checkPing Ping连通性检测
func (b *Monitor) checkPing(ctx context.Context, host string, timeout time.Duration, result *vmodel.HostCheckResult) bool {
	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pingCmd := fmt.Sprintf("ping -c 1 -W %d %s", monitorStatus.Config.CheckTimeout, host)
	b.XLog.Debugf("执行ping命令: %s", pingCmd)
	if output, err := utils.RunBashCtx(pingCtx, pingCmd); err != nil {
		result.Error = fmt.Sprintf("ping检测失败: %v", err)
		b.XLog.Debugf("ping命令失败 %s - 错误: %v, 输出: %s", host, err, string(output))
		return false
	} else {
		b.XLog.Debugf("ping命令成功 %s - 输出: %s", host, string(output))
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
	b.XLog.Errorf("网络连通性持续异常，执行路由器重启 - 连续失败 %d 次，达到重启阈值", 
		monitorStatus.Config.FailThreshold)

	// 更新重启统计
	now := time.Now()
	monitorStatus.LastRestartTime = &now
	monitorStatus.TotalRestarts++
	monitorStatus.RestartsInWindow++
	monitorStatus.ConsecutiveFails = 0
	monitorStatus.CurrentStatus = "cooldown"

	b.XLog.Infof("更新重启统计 - 总重启次数: %d, 窗口内重启次数: %d/%d, 冷却期: %d分钟",
		monitorStatus.TotalRestarts, monitorStatus.RestartsInWindow, 
		monitorStatus.Config.MaxRestarts, monitorStatus.Config.CooldownPeriod)

	// 记录重启日志
	reason := fmt.Sprintf("连续%d次网络检测失败", monitorStatus.Config.FailThreshold)
	if err := b.MonitorDao.LogRestart(reason); err != nil {
		b.XLog.Errorf("记录重启日志失败: %v", err)
	} else {
		b.XLog.Infof("重启日志已记录: %s", reason)
	}

	// 保存状态
	if err := b.MonitorDao.SaveMonitorStatus(monitorStatus); err != nil {
		b.XLog.Errorf("保存重启后状态失败: %v", err)
	}

	// 执行重启命令（异步）
	go func() {
		// 延迟执行，确保状态已保存
		b.XLog.Infof("延迟2秒后执行重启命令，确保状态保存完成")
		time.Sleep(2 * time.Second)
		
		b.XLog.Infof("开始执行系统重启命令: %s", cmds.ScriptReboot)
		if output, err := utils.RunBash(cmds.ScriptReboot); err != nil {
			b.XLog.Errorf("执行重启命令失败: %v, 输出: %s", err, string(output))
		} else {
			b.XLog.Infof("路由器重启命令执行成功，输出: %s", string(output))
		}
	}()
}

// calculateAverageLatency 计算平均延迟（毫秒）
func (b *Monitor) calculateAverageLatency(results []vmodel.HostCheckResult) float64 {
	var totalLatency time.Duration
	successCount := 0
	
	for _, result := range results {
		if result.Success {
			totalLatency += result.Latency
			successCount++
		}
	}
	
	if successCount == 0 {
		return 0
	}
	
	avgNanos := float64(totalLatency.Nanoseconds()) / float64(successCount)
	return avgNanos / 1e6 // 转换为毫秒
}

// ========== 定时任务相关方法 ==========

// NetworkMonitorJob 网络监控定时任务执行函数
func (b *Monitor) NetworkMonitorJob(ctx context.Context) {
	b.XLog.Infof("定时任务触发 - 开始执行网络监控检测 [%s]", time.Now().Format("2006-01-02 15:04:05"))
	
	// 记录执行开始时间
	startTime := time.Now()
	
	// 调用网络检测逻辑
	b.PerformNetworkCheck(ctx)
	
	// 记录执行完成时间和耗时
	duration := time.Since(startTime)
	b.XLog.Infof("网络监控检测执行完成 - 总耗时: %.0fms [%s]", 
		float64(duration.Nanoseconds())/1e6, time.Now().Format("15:04:05"))
}

// InitNetworkMonitor 初始化网络监控（从配置文件读取配置）
func (b *Monitor) InitNetworkMonitorFromConfig() error {
	// 从配置文件获取网络监控配置，其他参数使用默认值
	config := &vmodel.NetworkMonitorConfig{
		Enable:            b.Config.NetworkMonitor.Enable,
		CheckInterval:     b.Config.NetworkMonitor.CheckInterval,
		FailCheckInterval: 60,  // 默认失败后检测间隔：60秒
		CheckTimeout:      10,  // 默认检测超时：10秒
		TestHosts:         b.Config.NetworkMonitor.TestHosts,
		FailThreshold:     3,   // 默认连续失败阈值：3次
		FailHostThreshold: 3,   // 默认失败主机数阈值：3个
		MaxRestarts:       5,   // 默认最大重启次数：5次
		RestartWindow:     24,  // 默认重启计数窗口：24小时
		CooldownPeriod:    30,  // 默认冷却期：30分钟
	}

	// 初始化监控状态
	if err := b.InitializeStatus(config); err != nil {
		b.XLog.Errorf("初始化网络监控状态失败: %v", err)
		return err
	}

	b.XLog.Info("网络监控初始化完成")
	return nil
}
