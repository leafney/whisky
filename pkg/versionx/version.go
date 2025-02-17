package versionx

// InfoSvc 版本信息结构体
type InfoSvc struct {
	Version   string
	GitCommit string
	BuildTime string
}

var (
	// VersionInfo 全局版本信息实例
	VersionInfo = &InfoSvc{}
)

// NewInfo 提供一个构造函数用于 wire 注入
func NewInfoSvc() *InfoSvc {
	return VersionInfo
}

// GetVersion 获取版本信息
func (v *InfoSvc) GetVersion() map[string]string {
	return map[string]string{
		"version":    VersionInfo.Version,
		"git_commit": VersionInfo.GitCommit,
		"build_time": VersionInfo.BuildTime,
	}
}
