//go:build wireinject
// +build wireinject

/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 18:20
 * @Description:
 */

package cmd

import (
	"github.com/google/wire"
	"github.com/leafney/whisky/config"
	"github.com/leafney/whisky/internal"
	"github.com/leafney/whisky/pkg/xlogx"
)

type Injector struct {
	L *xlogx.XLogSvc
	R DefRouter
	C *config.Config
}

func BuildInjector(stop chan struct{}) (*Injector, func(), error) {
	wire.Build(
		internal.ProviderSet,
		wire.Struct(new(DefRouter), "*"),
		wire.Struct(new(Injector), "*"),
	)
	return &Injector{}, nil, nil
}
