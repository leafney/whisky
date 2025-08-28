package service

import (
	"context"

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
	// 网络监控任务现在由 NetworkTask Biz 层管理

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
	// 网络监控任务方法现在在 NetworkTask Biz 层中注册
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