// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/chainxx/bitx/internal/biz"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/chainxx/bitx/internal/data"
	"github.com/chainxx/bitx/internal/service"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confData *conf.Data, binance *conf.Binance, logger log.Logger) (*app, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	walletDbRepo := data.NewWalletDataSource(dataData)
	walletUseCase := biz.NewWalletUseCase(binance, walletDbRepo, logger)
	marketDbRepo := data.NewMarketDataResource(dataData)
	tradesUseCase := biz.NewTradesUseCase(binance, marketDbRepo, logger)
	lottoUseCase := biz.NewLottoUseCase()
	bitxService := service.NewMessageService(walletUseCase, tradesUseCase, lottoUseCase)
	mainApp := newApp(bitxService)
	return mainApp, func() {
		cleanup()
	}, nil
}
