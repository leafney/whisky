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
	"github.com/gofiber/fiber/v2"
	"github.com/leafney/whisky/internal/api"
)

type DefRouter struct {
	HomeApi    *api.Home
	RouterApi  *api.Router
	YacdApi    *api.YAcd
	NetWorkApi *api.NetWork
	SCrashApi  *api.SCrash
}

func (r *DefRouter) AutoMigrate() error {
	return nil
}

func (r *DefRouter) AutoStart(ctx context.Context) error {

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

func (r *DefRouter) UseMiddlewares(app *fiber.App, inj *Injector) {
	//app.Use(Cors())
}

func (r *DefRouter) SetupRoutes(app *fiber.App, inj *Injector) {

	app.Get("/", r.HomeApi.Home)
	app.Get("/version", r.HomeApi.Version)

	// router
	app.Get("/router", r.RouterApi.RouterInfo)
	app.Post("/router", r.RouterApi.RouterStatus)

	// network
	app.Get("/network", r.NetWorkApi.NetWorkInfo)

	//	clash
	app.Post("/scrash", r.SCrashApi.SCrashAction)
	//app.Post("/oclash", handler.YacdClashAction)

	app.Get("/yacd", r.YacdApi.YacdClashInfo)
	app.Post("/yacd", r.YacdApi.YacdClashAction)

	//// webui
	//uiDist, err := fs.Sub(web.UiStatic, "dist")
	//if err != nil {
	//	inj.L.Fatalf("static dir load error [%v]", err)
	//}
	//app.Use("/", filesystem.New(filesystem.Config{
	//	Root: http.FS(uiDist),
	//}))
}
