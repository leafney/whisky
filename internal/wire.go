/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 17:16
 * @Description:
 */

package internal

import (
	"github.com/google/wire"
	"github.com/leafney/whisky/internal/api"
	"github.com/leafney/whisky/internal/dao"
	"github.com/leafney/whisky/internal/service"
)

var Set = wire.NewSet(

	//
	wire.Struct(new(api.Home), "*"),
	wire.Struct(new(dao.Home), "*"),

	wire.Struct(new(service.NetWork), "*"),
	wire.Struct(new(service.Router), "*"),
	wire.Struct(new(service.YAcd), "*"),
	wire.Struct(new(service.Scrash), "*"),
)
