/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-05 23:31
 * @Description:
 */

package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/leafney/whisky/cmd"
	"github.com/leafney/whisky/pkg/versionx"
	"github.com/spf13/pflag"
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

func init() {
	// 初始化版本信息
	versionx.VersionInfo.Version = Version
	versionx.VersionInfo.GitCommit = GitCommit
	versionx.VersionInfo.BuildTime = BuildTime
}

func main() {
	pflag.BoolVarP(&h, "help", "h", false, "help")
	pflag.BoolVarP(&d, "debug", "d", false, "whether to output debug level logs")
	pflag.StringVarP(&p, "port", "p", "8080", "server port")
	pflag.StringVarP(&y, "yacd", "y", "9999", "yacd port")
	pflag.StringVarP(&w, "webhook", "w", "", "webhook url")
	pflag.BoolVarP(&v, "version", "v", false, "version")
	pflag.Parse()

	if h {
		pflag.PrintDefaults()
	} else if v {
		// 输出版本信息
		fmt.Println("Version:      " + Version)
		fmt.Println("Git branch:   " + GitBranch)
		fmt.Println("Git commit:   " + GitCommit)
		fmt.Println("Built time:   " + BuildTime)
		fmt.Println("Go version:   " + runtime.Version())
		fmt.Println("OS/Arch:      " + runtime.GOOS + "/" + runtime.GOARCH)
	} else {
		//// 基础服务
		//core.InitXLog(d)
		////core.InitConfig()
		//core.InitEConfig(y, w)
		////core.InitMsqQueue()
		//core.InitShellClean()
		//
		//// 用于退出的通道
		//quitChan := make(chan struct{})
		//// 相关服务
		////core.InitLog(quitChan)
		////core.InitMongo(quitChan)
		////core.InitRedis(quitChan)
		////core.InitCron(quitChan)
		////core.InitCache(quitChan)
		////core.InitNotify()
		//core.InitLevelDB(quitChan)
		//// 异步任务
		////core.InitQueue()
		//// web服务
		//run.Start(p, quitChan)

		// 用于退出的通道
		quitChan := make(chan struct{})
		injector, callback, err := cmd.BuildInjector(quitChan)
		if err != nil {
			log.Fatalln(err)
		}
		defer callback()

		ctx := context.Background()
		if err := injector.R.Init(ctx); err != nil {
			log.Fatalf("初始化异常 %v", err)
		}

		cmd.StartServer(injector, p, quitChan)
	}
}
