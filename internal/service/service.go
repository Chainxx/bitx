package service

import (
	"github.com/chainxx/bitx/internal/biz"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewMessageService)

type BitxService struct {
	wallet *biz.WalletUseCase
	trades *biz.TradesUseCase
	lotto  *biz.LottoUseCase
}

func NewMessageService(wallet *biz.WalletUseCase, trades *biz.TradesUseCase, lotto *biz.LottoUseCase) *BitxService {
	return &BitxService{
		wallet: wallet,
		trades: trades,
		lotto:  lotto,
	}
}
