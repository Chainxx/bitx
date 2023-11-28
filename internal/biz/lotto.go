package biz

import (
	"context"
	"fmt"
	"github.com/chainxx/bitx/pkg/lotto"
)

type LottoUseCase struct {
}

func NewLottoUseCase() *LottoUseCase {
	return &LottoUseCase{}
}

type DoubleColorBallBetting struct {
	Count   int64  `json:"count,omitempty"`
	Cost    int64  `json:"cost,omitempty"`
	Message string `json:"message,omitempty"`
}

func (lottoUseCase *LottoUseCase) GetDoubleColorBallCount(ctx context.Context, redBallCount int, blueBallCount int) DoubleColorBallBetting {
	count := lotto.DoubleColorBallCount(redBallCount, blueBallCount)
	return DoubleColorBallBetting{
		Count: count,
		Cost:  count * 2,
		Message: fmt.Sprintf("您计划买入红球(%d)个,蓝球(%d),投注注数(%d),总花费(%d)",
			redBallCount, blueBallCount, count, count*2),
	}
}
