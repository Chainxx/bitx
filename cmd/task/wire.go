//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/chainxx/bitx/internal/biz"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/chainxx/bitx/internal/data"
	"github.com/chainxx/bitx/internal/task"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

func wireApp(*conf.Data, *conf.Binance, log.Logger) (*app, func(), error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, task.ProviderSet, newApp))
}
