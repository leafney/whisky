package cronx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/pkg/xlogx"
)

// TaskMethodInfo 定义任务方法类型和描述信息
type TaskMethodInfo struct {
	Method      func(ctx context.Context) `json:"-"`    // 方法统一为 func(ctx context.Context) 类型
	Description string                    `json:"desc"` // 描述信息
	Name        string                    `json:"name"` // 方法名称
}

// TaskJob 存储任务信息
type TaskJob struct {
	Job    *gocron.Job
	Method func(ctx context.Context)
}

type CronSvc struct {
	scheduler   *gocron.Scheduler
	tasks       map[string]*TaskJob // 修改为存储 TaskJob
	taskMethods map[string]TaskMethodInfo
	ctx         context.Context
}

func NewCronSvc(cfg *config.Config, log *xlogx.XLogSvc, stop chan struct{}) *CronSvc {
	// 需要注意启动时的时区问题
	scheduler := gocron.NewScheduler(time.Local)

	go func() {
		<-stop
		scheduler.Stop()
		log.Infoln("[Cron] Exit successful")
	}()

	// 启动定时任务，立即启动或者在外部按需启动
	//scheduler.StartAsync()

	log.Infoln("[Cron] Load successful")

	return &CronSvc{
		scheduler:   scheduler,
		tasks:       make(map[string]*TaskJob),
		taskMethods: make(map[string]TaskMethodInfo),
		ctx:         context.Background(),
	}
}

// RegisterTaskMethod 注册任务方法
func (c *CronSvc) RegisterTaskMethod(name string, desc string, method func(ctx context.Context)) {
	if c.taskMethods == nil {
		c.taskMethods = make(map[string]TaskMethodInfo)
	}
	c.taskMethods[name] = TaskMethodInfo{
		Method:      method,
		Description: desc,
		Name:        name,
	}
}

// GetRegisteredMethods 获取所有注册的任务方法
func (c *CronSvc) GetRegisteredMethods() map[string]TaskMethodInfo {
	return c.taskMethods
}

// GetRegisteredMethodsWithDesc 获取所有注册的任务方法名称和描述
func (c *CronSvc) GetRegisteredMethodsWithDesc() []TaskMethodInfo {
	methods := make([]TaskMethodInfo, 0, len(c.taskMethods))
	for name, info := range c.taskMethods {
		methods = append(methods, TaskMethodInfo{
			Name:        name,
			Description: info.Description,
		})
	}
	return methods
}

// AddJob 添加定时任务
func (c *CronSvc) AddJob(taskId string, cronExpr string, taskFunc func(ctx context.Context)) error {
	// 检查任务是否已经存在
	if _, exists := c.tasks[taskId]; exists {
		return fmt.Errorf("task with ID [%s] already exists", taskId)
	}

	// 添加定时任务
	job, err := c.scheduler.Cron(cronExpr).Tag(taskId).SingletonMode().Do(taskFunc, c.ctx)
	if err != nil {
		return err
	}

	// 存储任务
	c.tasks[taskId] = &TaskJob{
		Job:    job,
		Method: taskFunc,
	}
	return nil
}

// AddJobSecs 添加定时任务，支持秒级任务
func (c *CronSvc) AddJobSecs(taskId string, cronExpr string, taskFunc func(ctx context.Context)) error {
	// 检查任务是否已经存在
	if _, exists := c.tasks[taskId]; exists {
		return fmt.Errorf("task with ID [%s] already exists", taskId)
	}

	// 添加定时任务
	job, err := c.scheduler.CronWithSeconds(cronExpr).Tag(taskId).SingletonMode().Do(taskFunc, c.ctx)
	if err != nil {
		return err
	}

	// 存储任务
	c.tasks[taskId] = &TaskJob{
		Job:    job,
		Method: taskFunc,
	}
	return nil
}

// UpdateTask 动态更新任务
func (c *CronSvc) UpdateJob(taskId string, cronExpr string, taskFunc func(ctx context.Context)) error {
	// 先删除旧任务
	if err := c.RemoveJob(taskId); err != nil {
		return err
	}

	// 添加新任务
	return c.AddJob(taskId, cronExpr, taskFunc)
}

// RemoveTask 动态删除任务
func (c *CronSvc) RemoveJob(taskId string) error {
	// 检查任务是否存在
	taskJob, exists := c.tasks[taskId]
	if !exists {
		return fmt.Errorf("task with ID [%s] does not exist", taskId)
	}

	// 删除任务
	c.scheduler.RemoveByReference(taskJob.Job)
	delete(c.tasks, taskId)
	return nil
}

// 判断指定任务是否存在
func (c *CronSvc) IsJobExists(taskId string) bool {
	_, exists := c.tasks[taskId]
	return exists
}

// Start 启动调度器
func (c *CronSvc) Start() {
	c.scheduler.StartAsync()
}

// Stop 停止调度器
func (c *CronSvc) Stop() {
	c.scheduler.Stop()
}

// GetJob 获取定时任务对象
func (c *CronSvc) GetJob(taskId string) *gocron.Job {
	if taskJob, exists := c.tasks[taskId]; exists {
		return taskJob.Job
	}
	return nil
}

// RunJobNow 立即执行一次任务
func (c *CronSvc) RunJobNow(taskId string) error {
	taskJob, exists := c.tasks[taskId]
	if !exists {
		return fmt.Errorf("task with ID [%s] does not exist", taskId)
	}

	// 使用 goroutine 异步执行任务，不影响调用方法的返回
	go taskJob.Method(c.ctx)
	return nil
}
