/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-08-13
 * @Description: 网络监控数据访问层
 */

package dao

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/leafney/whisky/internal/vmodel"
	"github.com/leafney/whisky/pkg/leveldbx"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Monitor struct {
	XLog    *xlogx.XLogSvc
	LevelDB *leveldbx.LevelDBSvc
}

const (
	// 存储键前缀
	keyMonitorStatus = "network_monitor:status"
	keyMonitorStats  = "network_monitor:stats"
	keyMonitorConfig = "network_monitor:config"
	keyRestartLog    = "network_monitor:restart_log:%s" // 格式: restart_log:2025-08-13
)

// SaveMonitorStatus 保存监控状态
func (d *Monitor) SaveMonitorStatus(status *vmodel.NetworkMonitorStatus) error {
	data, err := json.Marshal(status)
	if err != nil {
		d.XLog.Errorf("序列化监控状态失败: %v", err)
		return err
	}

	if err := d.LevelDB.SetS(keyMonitorStatus, string(data)); err != nil {
		d.XLog.Errorf("保存监控状态失败: %v", err)
		return err
	}

	d.XLog.Debug("监控状态已保存")
	return nil
}

// GetMonitorStatus 获取监控状态
func (d *Monitor) GetMonitorStatus() (*vmodel.NetworkMonitorStatus, error) {
	data, err := d.LevelDB.GetS(keyMonitorStatus)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			// 返回默认状态
			return &vmodel.NetworkMonitorStatus{
				Enabled:       false,
				CurrentStatus: "disabled",
			}, nil
		}
		d.XLog.Errorf("获取监控状态失败: %v", err)
		return nil, err
	}

	var status vmodel.NetworkMonitorStatus
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		d.XLog.Errorf("反序列化监控状态失败: %v", err)
		return nil, err
	}

	return &status, nil
}

// SaveMonitorConfig 保存监控配置
func (d *Monitor) SaveMonitorConfig(config *vmodel.NetworkMonitorConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		d.XLog.Errorf("序列化监控配置失败: %v", err)
		return err
	}

	if err := d.LevelDB.SetS(keyMonitorConfig, string(data)); err != nil {
		d.XLog.Errorf("保存监控配置失败: %v", err)
		return err
	}

	d.XLog.Debug("监控配置已保存")
	return nil
}

// GetMonitorConfig 获取监控配置
func (d *Monitor) GetMonitorConfig() (*vmodel.NetworkMonitorConfig, error) {
	data, err := d.LevelDB.GetS(keyMonitorConfig)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			return nil, nil // 配置不存在
		}
		d.XLog.Errorf("获取监控配置失败: %v", err)
		return nil, err
	}

	var config vmodel.NetworkMonitorConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		d.XLog.Errorf("反序列化监控配置失败: %v", err)
		return nil, err
	}

	return &config, nil
}

// SaveMonitorStats 保存监控统计信息
func (d *Monitor) SaveMonitorStats(stats *vmodel.NetworkMonitorStats) error {
	data, err := json.Marshal(stats)
	if err != nil {
		d.XLog.Errorf("序列化监控统计失败: %v", err)
		return err
	}

	if err := d.LevelDB.SetS(keyMonitorStats, string(data)); err != nil {
		d.XLog.Errorf("保存监控统计失败: %v", err)
		return err
	}

	return nil
}

// GetMonitorStats 获取监控统计信息
func (d *Monitor) GetMonitorStats() (*vmodel.NetworkMonitorStats, error) {
	data, err := d.LevelDB.GetS(keyMonitorStats)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			// 返回默认统计
			return &vmodel.NetworkMonitorStats{
				TotalChecks:      0,
				SuccessfulChecks: 0,
				FailedChecks:     0,
				SuccessRate:      0.0,
			}, nil
		}
		d.XLog.Errorf("获取监控统计失败: %v", err)
		return nil, err
	}

	var stats vmodel.NetworkMonitorStats
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		d.XLog.Errorf("反序列化监控统计失败: %v", err)
		return nil, err
	}

	return &stats, nil
}

// LogRestart 记录重启日志
func (d *Monitor) LogRestart(reason string) error {
	now := time.Now()
	dateKey := fmt.Sprintf(keyRestartLog, now.Format("2006-01-02"))

	// 获取当天已有的重启记录
	var restartLogs []map[string]interface{}
	data, err := d.LevelDB.GetS(dateKey)
	if err == nil {
		json.Unmarshal([]byte(data), &restartLogs)
	}

	// 添加新的重启记录
	logEntry := map[string]interface{}{
		"time":   now,
		"reason": reason,
	}
	restartLogs = append(restartLogs, logEntry)

	// 保存更新后的记录
	newData, err := json.Marshal(restartLogs)
	if err != nil {
		d.XLog.Errorf("序列化重启日志失败: %v", err)
		return err
	}

	if err := d.LevelDB.SetS(dateKey, string(newData)); err != nil {
		d.XLog.Errorf("保存重启日志失败: %v", err)
		return err
	}

	d.XLog.Infof("重启日志已记录: %s", reason)
	return nil
}

// GetRestartLogs 获取指定日期范围内的重启日志
func (d *Monitor) GetRestartLogs(days int) ([]map[string]interface{}, error) {
	var allLogs []map[string]interface{}

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dateKey := fmt.Sprintf(keyRestartLog, date)

		data, err := d.LevelDB.GetS(dateKey)
		if err != nil {
			continue // 该日期没有记录
		}

		var dayLogs []map[string]interface{}
		if err := json.Unmarshal([]byte(data), &dayLogs); err != nil {
			d.XLog.Errorf("反序列化重启日志失败 [%s]: %v", date, err)
			continue
		}

		allLogs = append(allLogs, dayLogs...)
	}

	return allLogs, nil
}

// ClearOldRestartLogs 清理旧的重启日志（保留指定天数）
func (d *Monitor) ClearOldRestartLogs(keepDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -keepDays)

	// 这里简化处理，实际可以遍历所有键进行清理
	for i := keepDays; i < keepDays+30; i++ { // 清理额外30天的旧数据
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dateKey := fmt.Sprintf(keyRestartLog, date)

		if err := d.LevelDB.Del(dateKey); err != nil {
			// 删除失败不影响主流程，只记录日志
			d.XLog.Debugf("删除旧重启日志失败 [%s]: %v", date, err)
		}
	}

	d.XLog.Infof("清理了 %s 之前的重启日志", cutoffDate.Format("2006-01-02"))
	return nil
}
