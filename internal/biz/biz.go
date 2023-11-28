package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTradesUseCase, NewWalletUseCase, NewLottoUseCase)
