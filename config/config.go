/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-09 10:23
 * @Description:
 */

package config

type Config struct {
	YacdPort string `json:"yacd_port"`
	WebHook  string `json:"web_hook"`
	Log      Log
	LevelDB  LevelDB
}

type (
	Log struct {
		XEnable bool `koanf:"xenable" default:"true"`
		XDebug  bool `default:"true"`
		XLevel  string
	}

	LevelDB struct {
		Path string `default:"data/.cache"`
	}
)
