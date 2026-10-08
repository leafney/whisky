/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:41
 * @Description:
 */

package xlogx

import (
	"github.com/leafney/rose/xlog"
	"github.com/leafney/whisky/config"
)

type XLogSvc struct {
	*xlog.Log
}

func NewXLogSvc(cfg *config.Config) *XLogSvc {
	cfgLog := cfg.Log

	x := xlog.NewXLog()
	// debug model
	x.SetDebug(cfgLog.XDebug)
	// enable
	x.SetEnable(cfgLog.XEnable)
	// default log level
	x.SetLevelStr(cfgLog.XLevel)

	x.Infoln("[XLog] Load successful")

	return &XLogSvc{x}
}
