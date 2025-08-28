/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:44
 * @Description:
 */

package cmd

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/leafney/whisky/internal/api"
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/web"
)

type DefRouter struct {
	HomeApi    *api.Home
	RouterApi  *api.Router
	YacdApi    *api.YAcd
	NetWorkApi *api.NetWork
	SCrashApi  *api.SCrash

	CronSvc        *service.Cron  // service 层的定时任务管理
	CronTaskApi    *api.CronTask
	NetworkTaskBiz *biz.NetworkTask
}

func (r *DefRouter) AutoMigrate() error {
	return nil
}

func (r *DefRouter) AutoStart(ctx context.Context) error {
	// 启动定时任务服务
	if r.CronSvc != nil {
		// 加载任务（包括 TestJob 等）
		if err := r.CronSvc.LoadJobs(ctx); err != nil {
			return err
		}

		// 启动调度器
		r.CronSvc.Start()

		// 自动启动网络监控（如果配置启用）
		if r.NetworkTaskBiz != nil {
			if err := r.autoStartNetworkMonitor(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *DefRouter) Init(ctx context.Context) error {

	if err := r.AutoMigrate(); err != nil {
		return err
	}
	if err := r.AutoStart(ctx); err != nil {
		return err
	}

	return nil
}

// Shutdown 优雅关闭
func (r *DefRouter) Shutdown() error {
	// 关闭定时任务服务
	if r.CronSvc != nil {
		r.CronSvc.Stop()
	}

	return nil
}

// autoStartNetworkMonitor 自动启动网络监控任务
func (r *DefRouter) autoStartNetworkMonitor() error {
	// 检查功能是否启用
	featureStatus := r.NetworkTaskBiz.GetTaskFeatureStatus()
	if enabled, ok := featureStatus["feature_enabled"].(bool); !ok || !enabled {
		return nil // 功能未启用，跳过
	}

	// 启动网络监控任务
	if err := r.NetworkTaskBiz.StartTask(); err != nil {
		return err
	}

	return nil
}

func (r *DefRouter) UseMiddlewares(app *fiber.App, inj *Injector) {
	//app.Use(Cors())
}

func (r *DefRouter) SetupRoutes(app *fiber.App, inj *Injector) {

	app.Get("/", r.HomeApi.Home)
	app.Get("/version", r.HomeApi.Version)

	// router
	app.Get("/router", r.RouterApi.RouterInfo)
	app.Post("/router", r.RouterApi.RouterStatus)

	// cron tasks - 定时任务管理
	app.Get("/cron/tasks", r.CronTaskApi.GetAllTasksStatus)
	app.Get("/cron/methods", r.CronTaskApi.GetTaskMethods)
	app.Post("/cron/task/:taskId", r.CronTaskApi.ControlTask)
	app.Get("/cron/task/:taskId/status", r.CronTaskApi.GetTaskStatus)

	// network
	app.Get("/network", r.NetWorkApi.NetWorkInfo)

	//	clash
	app.Post("/scrash", r.SCrashApi.SCrashAction)
	//app.Post("/oclash", handler.YacdClashAction)

	app.Get("/yacd", r.YacdApi.YacdClashInfo)
	app.Post("/yacd", r.YacdApi.YacdClashAction)

	// webui
	uiDist, err := web.GetDistFS()
	if err != nil {
		inj.L.Fatalf("static dir load error [%v]", err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(uiDist),
	}))

	// 处理所有未匹配的请求，返回 index.html
	app.Get("/*", func(c *fiber.Ctx) error {
		// 从嵌入的文件系统中打开 index.html
		indexFile, err := web.GetDistFSRoot().Open("dist/index.html")
		if err != nil {
			println("打开 index.html 失败:", err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("无法加载 index.html")
		}
		defer indexFile.Close()

		// 设置 MIME 类型并发送文件流
		c.Type("html")
		return c.SendStream(indexFile)
	})
}
