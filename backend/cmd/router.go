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
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/web"
)

type DefRouter struct {
	HomeApi           *api.Home
	RouterApi         *api.Router
	YacdApi           *api.YAcd
	NetWorkApi        *api.NetWork
	SCrashApi         *api.SCrash
	NetworkMonitorSvc *service.NetworkMonitor
}

func (r *DefRouter) AutoMigrate() error {
	return nil
}

func (r *DefRouter) AutoStart(ctx context.Context) error {
	// 初始化并自动启动网络监控（如果配置启用）
	if r.NetworkMonitorSvc != nil {
		// 先初始化调度器
		if err := r.NetworkMonitorSvc.Initialize(); err != nil {
			return err
		}

		// 自动启动监控
		if err := r.NetworkMonitorSvc.AutoStart(ctx); err != nil {
			return err
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
	// 关闭网络监控服务
	if r.NetworkMonitorSvc != nil {
		if err := r.NetworkMonitorSvc.Shutdown(); err != nil {
			return err
		}
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

	// network monitor
	app.Get("/router/monitor", r.RouterApi.NetworkMonitorStatus)
	app.Post("/router/monitor", r.RouterApi.NetworkMonitorControl)

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
