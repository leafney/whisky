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
	p         string
	v         bool
	h         bool
	Version   = "v0.1.0"
	GitBranch = ""
	GitCommit = ""
	BuildTime = "2024-07-06 13:04:30"
)

func main() {
	flag.BoolVarP(&h, "help", "h", false, "help")
	flag.StringVarP(&p, "port", "p", "8080", "port")
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
		core.InitXLog()
		//core.InitConfig()
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
