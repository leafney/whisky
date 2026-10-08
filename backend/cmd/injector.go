/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:34
 * @Description:
 */

package cmd

import (
	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Injector struct {
	L *xlogx.XLogSvc
	R DefRouter
	C *config.Config
}
