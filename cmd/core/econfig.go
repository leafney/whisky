/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-11 08:56
 * @Description:
 */

package core

import (
	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/global"
)

func InitEConfig(yacdPort, webHook string) {
	eCfg := &config.Config{
		YacdPort: yacdPort,
		WebHook:  webHook,
	}

	global.GEConfig = eCfg

	global.GXLog.Infoln("[EConfig] Load successful")
}
