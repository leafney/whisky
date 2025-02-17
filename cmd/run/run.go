/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 12:38
 * @Description:
 */

package run

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leafney/rose"
	"github.com/leafney/whisky/global"
)

func Start(port string, stop chan struct{}) {
	cfgPort := "8080"
	if !rose.StrIsEmpty(port) {
		cfgPort = port
	}

	app := fiber.New()

	//	router
	bindRouter(app)

	// shutdown
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
		close(stop)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			global.GXLog.Errorf("[Server] Shutdown error [%v]", err)
		} else {
			global.GXLog.Infoln("[Server] Shutdown successful")
		}
		time.Sleep(100 * time.Millisecond)
	}()

	// start
	global.GXLog.Infoln("[Server] Load successful")
	if err := app.Listen(fmt.Sprintf(":%s", cfgPort)); err != nil {
		global.GXLog.Errorf("[Server] Listen error [%v]", err)
	}

	wg.Wait()
	global.GXLog.Infoln("[Server] Exit successful")
}
