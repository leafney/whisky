/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 网络监控业务逻辑层（基于定时任务）
 */

package biz

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
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

	// 定时任务调度器
	scheduler gocron.Scheduler
	job       gocron.Job
	mutex     sync.RWMutex
	isRunning bool
}

// 全局监控状态（简化状态管理）
var (
	monitorStatus *vmodel.NetworkMonitorStatus
	monitorStats  *vmodel.NetworkMonitorStats
	stateMutex    sync.RWMutex
)

// InitScheduler 初始化调度器
func (b *Monitor) InitScheduler() error {
	var err error
	b.scheduler, err = gocron.NewScheduler()
	if err != nil {
		b.XLog.Errorf("创建调度器失败: %v", err)
		return err
	}

	b.XLog.Info("网络监控调度器初始化成功")
	return nil
}

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

// Start 启动网络监控定时任务
func (b *Monitor) Start(config *vmodel.NetworkMonitorConfig) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.isRunning {
		return fmt.Errorf("网络监控已在运行中")
	}

	// 初始化状态
	if err := b.initializeStatus(config); err != nil {
		return fmt.Errorf("初始化状态失败: %v", err)
	}

	// 创建定时任务
	checkInterval := time.Duration(config.CheckInterval) * time.Second
	job, err := b.scheduler.NewJob(
		gocron.DurationJob(checkInterval),
		gocron.NewTask(b.performNetworkCheck),
		gocron.WithName("network_monitor"),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	if err != nil {
		b.XLog.Errorf("创建网络监控定时任务失败: %v", err)
		return err
	}

	b.job = job
	b.isRunning = true

	// 更新状态
	stateMutex.Lock()
	monitorStatus.Enabled = true
	monitorStatus.CurrentStatus = "running"
	stateMutex.Unlock()

	// 启动调度器
	b.scheduler.Start()

	// 保存状态
	if err := b.MonitorDao.SaveMonitorStatus(monitorStatus); err != nil {
		b.XLog.Errorf("保存监控状态失败: %v", err)
	}

	b.XLog.Infof("网络监控定时任务已启动，检测间隔: %v", checkInterval)
	return nil
}

// Stop 停止网络监控定时任务
func (b *Monitor) Stop() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !b.isRunning {
		return fmt.Errorf("网络监控未在运行")
	}

	// 停止定时任务
	if b.job != nil {
		if err := b.scheduler.RemoveJob(b.job.ID()); err != nil {
			b.XLog.Errorf("移除定时任务失败: %v", err)
		}
		b.job = nil
	}

	b.isRunning = false

	// 更新状态
	stateMutex.Lock()
	if monitorStatus != nil {
		monitorStatus.Enabled = false
		monitorStatus.CurrentStatus = "disabled"
	}
	stateMutex.Unlock()

	// 保存状态
	if monitorStatus != nil {
		if err := b.MonitorDao.SaveMonitorStatus(monitorStatus); err != nil {
			b.XLog.Errorf("保存监控状态失败: %v", err)
		}
	}

	b.XLog.Info("网络监控定时任务已停止")
	return nil
}

// Restart 重启网络监控
func (b *Monitor) Restart(config *vmodel.NetworkMonitorConfig) error {
	if err := b.Stop(); err != nil {
		b.XLog.Errorf("停止监控失败: %v", err)
	}

	// 等待一秒确保完全停止
	time.Sleep(1 * time.Second)

	return b.Start(config)
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

// UpdateConfig 更新配置并重新调度任务
func (b *Monitor) UpdateConfig(config *vmodel.NetworkMonitorConfig) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	stateMutex.Lock()
	if monitorStatus != nil {
		monitorStatus.Config = *config
	}
	stateMutex.Unlock()

	// 保存配置
	if err := b.MonitorDao.SaveMonitorConfig(config); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	if err := b.MonitorDao.SaveMonitorStatus(monitorStatus); err != nil {
		return fmt.Errorf("保存状态失败: %v", err)
	}

	// 如果正在运行，重新创建任务以应用新配置
	if b.isRunning {
		b.XLog.Info("配置已更新，重新启动监控任务")
		return b.Restart(config)
	}

	b.XLog.Info("网络监控配置已更新")
	return nil
}

// Shutdown 优雅关闭调度器
func (b *Monitor) Shutdown() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.scheduler != nil {
		if err := b.scheduler.Shutdown(); err != nil {
			b.XLog.Errorf("关闭调度器失败: %v", err)
			return err
		}
		b.XLog.Info("网络监控调度器已关闭")
	}

	return nil
}

// 初始化状态
func (b *Monitor) initializeStatus(config *vmodel.NetworkMonitorConfig) error {
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

// performNetworkCheck 执行网络检测（定时任务回调函数）
func (b *Monitor) performNetworkCheck() {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	// 检查是否在冷却期
	if b.isInCooldown() {
		b.XLog.Debug("处于冷却期，跳过检测")
		return
	}

	b.XLog.Debug("开始执行网络连通性检测")
	now := time.Now()
	monitorStatus.LastCheckTime = &now

	// 并发检测所有主机
	results := b.checkAllHosts()
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
		} else {
			// 失败时调整检测频率
			b.adjustCheckFrequency(true)
		}
	} else {
		// 检测成功，重置失败计数
		if monitorStatus.ConsecutiveFails > 0 {
			b.XLog.Info("网络连通性恢复正常")
			monitorStatus.ConsecutiveFails = 0
			// 恢复正常检测频率
			b.adjustCheckFrequency(false)
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

// adjustCheckFrequency 调整检测频率
func (b *Monitor) adjustCheckFrequency(isFailed bool) {
	if !b.isRunning || b.job == nil {
		return
	}

	var newInterval time.Duration
	if isFailed {
		newInterval = time.Duration(monitorStatus.Config.FailCheckInterval) * time.Second
		b.XLog.Debugf("调整为失败检测间隔: %v", newInterval)
	} else {
		newInterval = time.Duration(monitorStatus.Config.CheckInterval) * time.Second
		b.XLog.Debugf("恢复正常检测间隔: %v", newInterval)
	}

	// 重新创建任务以更新间隔
	go func() {
		b.mutex.Lock()
		defer b.mutex.Unlock()

		if b.job != nil {
			// 移除旧任务
			if err := b.scheduler.RemoveJob(b.job.ID()); err != nil {
				b.XLog.Errorf("移除旧任务失败: %v", err)
				return
			}

			// 创建新任务
			newJob, err := b.scheduler.NewJob(
				gocron.DurationJob(newInterval),
				gocron.NewTask(b.performNetworkCheck),
				gocron.WithName("network_monitor"),
			)

			if err != nil {
				b.XLog.Errorf("创建新任务失败: %v", err)
				return
			}

			b.job = newJob
		}
	}()
}

// 检测所有主机
func (b *Monitor) checkAllHosts() []vmodel.HostCheckResult {
	hosts := monitorStatus.Config.TestHosts
	results := make([]vmodel.HostCheckResult, len(hosts))
	var wg sync.WaitGroup

	for i, host := range hosts {
		wg.Add(1)
		go func(index int, hostname string) {
			defer wg.Done()
			results[index] = b.pingHost(hostname)
		}(i, host)
	}

	wg.Wait()
	return results
}

// Ping单个主机
func (b *Monitor) pingHost(host string) vmodel.HostCheckResult {
	start := time.Now()
	result := vmodel.HostCheckResult{
		Host:      host,
		CheckTime: start,
	}

	// 设置超时
	timeout := time.Duration(monitorStatus.Config.CheckTimeout) * time.Second

	// 先尝试TCP连接
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), timeout)
	if err == nil {
		conn.Close()
		result.Success = true
		result.Latency = time.Since(start)
		return result
	}

	// TCP连接失败，尝试ping
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pingCmd := fmt.Sprintf("ping -c 1 -W %d %s", monitorStatus.Config.CheckTimeout, host)
	if _, pingErr := utils.RunBashCtx(ctx, pingCmd); pingErr != nil {
		result.Success = false
		result.Error = fmt.Sprintf("TCP连接和ping都失败: %v, %v", err, pingErr)
	} else {
		result.Success = true
		result.Latency = time.Since(start)
	}

	return result
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
		// 停止定时任务
		go func() {
			time.Sleep(1 * time.Second)
			b.Stop()
		}()
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
