/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2025-02-17 18:20
 * @Description:
 */

package cmd

import "github.com/google/wire"

func BuildInjector(stop chan struct{}) (*Injector, func(), error) {
	wire.Build(
		AppSet,
		wire.Struct(new(DefRouter), "*"),
		wire.Struct(new(Injector), "*"),
	)
	return &Injector{}, nil, nil
}
