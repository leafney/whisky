package versionx

// InfoSvc 版本信息结构体
type InfoSvc struct {
	Version   string
	GitCommit string
	BuildTime string
}

// VersionInfo 全局版本信息实例
var VersionInfo = &InfoSvc{}

// NewInfo 提供一个构造函数用于 wire 注入
func NewInfoSvc() *InfoSvc {
	return VersionInfo
}

// SetVersion 设置版本信息
func SetVersion(version, commit, buildTime string) {
	VersionInfo.Version = version
	VersionInfo.GitCommit = commit
	VersionInfo.BuildTime = buildTime
}

// GetVersion 获取版本信息
func (v *InfoSvc) GetVersion() map[string]string {
	return map[string]string{
		"version":    v.Version,
		"git_commit": v.GitCommit,
		"build_time": v.BuildTime,
	}
}
