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
	//
	wire.Struct(new(api.NetWork), "*"),
	// wire.Struct(new(dao.NetWork), "*"),
	//
	wire.Struct(new(api.Router), "*"),
	// wire.Struct(new(dao.Router), "*"),
	//
	wire.Struct(new(api.OClash), "*"),
	// wire.Struct(new(dao.OClash), "*"),
	//
	wire.Struct(new(api.SCrash), "*"),
	// wire.Struct(new(dao.SCrash), "*"),
	//
	wire.Struct(new(api.YAcd), "*"),
	// wire.Struct(new(dao.YAcd), "*"),

	wire.Struct(new(service.NetWork), "*"),
	wire.Struct(new(service.Router), "*"),
	wire.Struct(new(service.YAcd), "*"),
	wire.Struct(new(service.SCrash), "*"),
)
