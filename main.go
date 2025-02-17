/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-05 23:31
 * @Description:
 */

package main

import (
	"fmt"
	"github.com/leafney/whisky/cmd/core"
	"github.com/leafney/whisky/cmd/run"
	flag "github.com/spf13/pflag"
	"runtime"
)

var (
	v         bool
	h         bool
	d         bool
	p         string
	y         string
	w         string
	Version   = "v0.1.0"
	GitBranch = ""
	GitCommit = ""
	BuildTime = "2024-07-06 13:04:30"
)

func main() {
	flag.BoolVarP(&h, "help", "h", false, "help")
	flag.BoolVarP(&h, "debug", "d", false, "whether to output debug level logs")
	flag.StringVarP(&p, "port", "p", "8080", "server port")
	flag.StringVarP(&y, "yacd", "y", "9999", "yacd port")
	flag.StringVarP(&w, "webhook", "w", "", "webhook url")
	flag.BoolVarP(&v, "version", "v", false, "version")
	flag.Parse()

	if h {
		flag.PrintDefaults()
	} else if v {
		// 输出版本信息
		fmt.Println("Version:      " + Version)
		fmt.Println("Git branch:   " + GitBranch)
		fmt.Println("Git commit:   " + GitCommit)
		fmt.Println("Built time:   " + BuildTime)
		fmt.Println("Go version:   " + runtime.Version())
		fmt.Println("OS/Arch:      " + runtime.GOOS + "/" + runtime.GOARCH)
	} else {
		// 基础服务
		core.InitXLog(d)
		//core.InitConfig()
		core.InitEConfig(y, w)
		//core.InitMsqQueue()
		core.InitShellClean()

		// 用于退出的通道
		quitChan := make(chan struct{})
		// 相关服务
		//core.InitLog(quitChan)
		//core.InitMongo(quitChan)
		//core.InitRedis(quitChan)
		//core.InitCron(quitChan)
		//core.InitCache(quitChan)
		//core.InitNotify()
		core.InitLevelDB(quitChan)
		// 异步任务
		//core.InitQueue()
		// web服务
		run.Start(p, quitChan)
	}
}
