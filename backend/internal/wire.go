/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:16
 * @Description: Wire 依赖注入配置
 */

package internal

import (
	"github.com/google/wire"
	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/internal/api"
	"github.com/leafney/whisky/internal/biz"
	"github.com/leafney/whisky/internal/dao"
	"github.com/leafney/whisky/internal/service"
	"github.com/leafney/whisky/pkg/cronx"
	"github.com/leafney/whisky/pkg/leveldbx"
	"github.com/leafney/whisky/pkg/versionx"
	"github.com/leafney/whisky/pkg/xlogx"
)

// ProviderSet 统一的依赖注入配置
var ProviderSet = wire.NewSet(
	// 基础设施层
	config.NewConfig,
	xlogx.NewXLogSvc,
	leveldbx.NewLevelDBSvc,
	versionx.NewInfoSvc,
	cronx.NewCronSvc,

	// 数据访问层
	wire.Struct(new(dao.Monitor), "*"),
	wire.Struct(new(dao.Home), "*"),

	// 业务逻辑层
	wire.Struct(new(biz.Monitor), "*"),
	wire.Struct(new(biz.Home), "*"),

	// 服务层（内部逻辑）
	wire.Struct(new(service.NetworkMonitor), "*"),
	wire.Struct(new(service.Cron), "*"),
	wire.Struct(new(service.NetWork), "*"),
	wire.Struct(new(service.Router), "*"),
	wire.Struct(new(service.YAcd), "*"),
	wire.Struct(new(service.SCrash), "*"),

	// API层（对外接口）
	wire.Struct(new(api.Home), "*"),
	wire.Struct(new(api.NetWork), "*"),
	wire.Struct(new(api.Router), "*"),
	wire.Struct(new(api.OClash), "*"),
	wire.Struct(new(api.SCrash), "*"),
	wire.Struct(new(api.YAcd), "*"),
	wire.Struct(new(api.CronTask), "*"),
)
