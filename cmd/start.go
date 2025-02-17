/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 18:14
 * @Description:
 */

package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func StartServer(injector *Injector, port string, quit chan struct{}) {

	f := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// default middlewares
	//f.Use()

	// middlewares
	injector.R.UseMiddlewares(f, injector)

	// routers
	injector.R.SetupRoutes(f, injector)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// 添加信号监听，支持优雅退出
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		// 等待停止信号
		<-signalChan

		// 关闭通道
		close(quit)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := f.ShutdownWithContext(ctx); err != nil {
			injector.L.Errorf("[Server] Shutdown error [%v]", err)
		} else {
			injector.L.Info("[Server] Shutdown successful")
		}
		time.Sleep(200 * time.Millisecond)
	}()

	// start
	injector.L.Info("[Server] Load successful")

	// 监听端口，默认 > 配置文件 > 命令行
	defPort := "8080"
	//defPort := vars.WebDefaultPort
	//if !rose.StrIsEmpty(injector.C.Port) {
	//	defPort = injector.C.Port
	//}
	//if !rose.StrIsEmpty(port) && port != vars.WebDefaultPort {
	//	defPort = port
	//}

	if err := f.Listen(fmt.Sprintf(":%s", defPort)); err != nil {
		injector.L.Fatalf("[Server] Listen error [%v]", err)
	}

	wg.Wait()
	injector.L.Info("[Server] Exit successful")

}
