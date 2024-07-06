/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:33
 * @Description:
 */

package core

import (
	"github.com/leafney/rose/xlog"
	"github.com/leafney/whisky/global"
)

func InitXLog() {

	x := xlog.NewXLog()

	//// TODO 设置启用开关
	//showLog := os.Getenv("DEBUG_LOG")
	//if rose.StrToLower(showLog) == "true" {
	//	cfg.SetEnable(true)
	//} else {
	//	cfg.SetEnable(false)
	//}

	// 调试
	x.SetDebug(true)
	// 日志级别
	//x.SetLevel(xlog.ErrorLevel)

	// 启用开关
	//x.SetEnable(false)

	global.GXLog = x

	global.GXLog.Infoln("[XLog] Load successful")
}
