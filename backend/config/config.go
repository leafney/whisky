/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-09 10:23
 * @Description:
 */

package config

import (
	"embed"
	"io/fs"
	"log"
	"os"

	"github.com/creasty/defaults"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/leafney/rose"
)

//go:embed config.toml.default
var DefaultConfig embed.FS

type Config struct {
	// YacdPort string `json:"yacd_port"`
	// WebHook  string `json:"web_hook"`
	Port           string `koanf:"port" default:"8080"`
	Log            Log
	LevelDB        LevelDB
	NetworkMonitor NetworkMonitor
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

	NetworkMonitor struct {
		Enable            bool     `koanf:"enable" default:"false"`
		CheckInterval     int      `koanf:"check_interval" default:"300"`
		FailCheckInterval int      `koanf:"fail_check_interval" default:"60"`
		CheckTimeout      int      `koanf:"check_timeout" default:"10"`
		TestHosts         []string `koanf:"test_hosts" default:"[\"8.8.8.8\",\"114.114.114.114\",\"1.1.1.1\",\"223.5.5.5\",\"www.baidu.com\"]"`
		FailThreshold     int      `koanf:"fail_threshold" default:"3"`
		FailHostThreshold int      `koanf:"fail_host_threshold" default:"3"`
		MaxRestarts       int      `koanf:"max_restarts" default:"5"`
		RestartWindow     int      `koanf:"restart_window" default:"24"`
		CooldownPeriod    int      `koanf:"cooldown_period" default:"30"`
	}
)

var k = koanf.New(".")

func NewConfig() (*Config, error) {

	cfg := new(Config)

	path := "data/config.toml"

	if exist := loadDefaultConfig(path); !exist {
		log.Printf("[Koanf] Create default config file [%v]", path)
		os.Exit(0)
	}

	if err := k.Load(file.Provider(path), toml.Parser()); err != nil {
		log.Fatalf("[Koanf] Load config file [%v] error [%v]", path, err)
	}

	// 先设置默认值
	if err := defaults.Set(cfg); err != nil {
		log.Fatalf("[Koanf] Set default value error [%v]", err)
	}
	// 再解析配置参数
	if err := k.Unmarshal("", cfg); err != nil {
		log.Fatalf("[Koanf] Unmarshal config error [%v]", err)
	}

	log.Println("[Koanf] Load successful")

	log.Println(rose.JsonMarshalStr(cfg))

	return cfg, nil
}

func loadDefaultConfig(path string) (exist bool) {
	exist = true
	if !rose.FIsExist(path) {
		exist = false
		// 保证配置文件所在目录存在
		if err := rose.DirExistsEnsure(path); err != nil {
			log.Fatalf("[Koanf] Failed to create config directory: %v", err)
		}

		data, err := fs.ReadFile(DefaultConfig, "config.toml.default")
		if err != nil {
			log.Fatalf("[Koanf] Failed to read embedded config: %v", err)
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			log.Fatalf("[Koanf] Failed to write default config file: %v", err)
		}

		log.Println("[Koanf] Default config file created")
	}
	return exist
}
