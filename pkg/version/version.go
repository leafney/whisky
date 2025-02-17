package version

// Info 版本信息结构体
type Info struct {
	Version   string
	GitCommit string
	BuildTime string
}

var (
	// VersionInfo 全局版本信息实例
	VersionInfo = &Info{}
)

// NewInfo 提供一个构造函数用于 wire 注入
func NewInfo() *Info {
	return VersionInfo
}

// GetVersion 获取版本信息
func GetVersion() map[string]string {
	return map[string]string{
		"version":    VersionInfo.Version,
		"git_commit": VersionInfo.GitCommit,
		"build_time": VersionInfo.BuildTime,
	}
}
