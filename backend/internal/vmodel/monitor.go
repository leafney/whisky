/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 网络监控相关数据模型
 */

package vmodel

import "time"

// NetworkMonitorStatus 网络监控状态
type NetworkMonitorStatus struct {
	Enabled          bool                 `json:"enabled"`            // 是否启用
	CurrentStatus    string               `json:"current_status"`     // 当前状态: running, cooldown, disabled, error
	LastCheckTime    *time.Time           `json:"last_check_time"`    // 最后检测时间
	NextCheckTime    *time.Time           `json:"next_check_time"`    // 下次检测时间
	ConsecutiveFails int                  `json:"consecutive_fails"`  // 连续失败次数
	TotalRestarts    int                  `json:"total_restarts"`     // 总重启次数
	LastRestartTime  *time.Time           `json:"last_restart_time"`  // 最后重启时间
	RestartsInWindow int                  `json:"restarts_in_window"` // 窗口期内重启次数
	WindowStartTime  *time.Time           `json:"window_start_time"`  // 窗口期开始时间
	LastCheckResults []HostCheckResult    `json:"last_check_results"` // 最后一次检测结果
	Config           NetworkMonitorConfig `json:"config"`             // 当前配置
}

// HostCheckResult 主机检测结果
type HostCheckResult struct {
	Host      string        `json:"host"`       // 主机地址
	Success   bool          `json:"success"`    // 是否成功
	Latency   time.Duration `json:"latency"`    // 延迟时间
	Error     string        `json:"error"`      // 错误信息
	CheckTime time.Time     `json:"check_time"` // 检测时间
}

// NetworkMonitorConfig 网络监控配置
type NetworkMonitorConfig struct {
	Enable            bool     `json:"enable"`              // 是否启用
	CheckInterval     int      `json:"check_interval"`      // 正常检测间隔（秒）
	FailCheckInterval int      `json:"fail_check_interval"` // 失败后检测间隔（秒）
	CheckTimeout      int      `json:"check_timeout"`       // 单次检测超时（秒）
	TestHosts         []string `json:"test_hosts"`          // 测试主机列表
	FailThreshold     int      `json:"fail_threshold"`      // 连续失败阈值
	FailHostThreshold int      `json:"fail_host_threshold"` // 单次检测失败主机数阈值
	MaxRestarts       int      `json:"max_restarts"`        // 最大重启次数
	RestartWindow     int      `json:"restart_window"`      // 重启计数窗口期（小时）
	CooldownPeriod    int      `json:"cooldown_period"`     // 重启后冷却期（分钟）
}

// NetworkMonitorRequest 网络监控请求
type NetworkMonitorRequest struct {
	Action string `json:"action"` // start, stop, reset, run_now
}

// NetworkMonitorStats 网络监控统计信息
type NetworkMonitorStats struct {
	TotalChecks       int     `json:"total_checks"`        // 总检测次数
	SuccessfulChecks  int     `json:"successful_checks"`   // 成功检测次数
	FailedChecks      int     `json:"failed_checks"`       // 失败检测次数
	SuccessRate       float64 `json:"success_rate"`        // 成功率
	AverageLatency    int64   `json:"average_latency"`     // 平均延迟（毫秒）
	UptimePercentage  float64 `json:"uptime_percentage"`   // 在线时间百分比
	LastWeekRestarts  int     `json:"last_week_restarts"`  // 最近一周重启次数
	LastMonthRestarts int     `json:"last_month_restarts"` // 最近一月重启次数
}
