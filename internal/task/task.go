package task

import (
	"context"
	"github.com/chainxx/bitx/internal/biz"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/robfig/cron"
)

type Task struct {
	trades  *biz.TradesUseCase
	wallet  *biz.WalletUseCase
	symbols []string
	log     *log.Helper
}

func NewTask(trades *biz.TradesUseCase, wallet *biz.WalletUseCase, binance *conf.Binance, logger log.Logger) *Task {
	return &Task{
		trades:  trades,
		wallet:  wallet,
		symbols: binance.Symbols,
		log:     log.NewHelper(logger),
	}
}

var ProviderSet = wire.NewSet(NewTask)

func (task *Task) Cronjob() {
	c := cron.New()
	ctx := context.Background()
	
	//symbols := []string{"ETHUSDT", "DOGEUSDT", "LUNAUSDT", "BTCUSDT", "BNBUSDT"}
	symbols := task.symbols
	
	// 1. 一段时间内价格信息
	err := c.AddFunc("1 * * * * *", func() {
		for _, symbol := range symbols {
			task.trades.TickerPrice(ctx, symbol, 60)
		}
	})
	
	if err != nil {
		task.log.Errorf("配置添加价格信息 cronjob 失败, err: %s", err.Error())
	}
	
	// 2. 24hr 价格变动情况
	err = c.AddFunc("1 0 8 * * *", func() {
		for _, symbol := range symbols {
			task.trades.DayTicker(ctx, symbol)
		}
	})
	
	if err != nil {
		task.log.Errorf("配置添加价格变动 cronjob 失败, err: %s", err.Error())
	}
	
	// 3. 获取每日资产信息 - 每日8点获取
	err = c.AddFunc("1 59 7 * * *", func() {
		task.wallet.GetPerDayAssets(ctx)
	})
	
	if err != nil {
		task.log.Errorf("配置获取每日资产信息 cronjob 失败, err: %s", err.Error())
	}
	
	// 4. 获取每日充值数据 - 每日8点获取
	err = c.AddFunc("1 0 8 * * *", func() {
		task.wallet.RechargeHis(ctx)
	})
	
	if err != nil {
		task.log.Errorf("配置获取每日充值信息 cronjob 失败, err: %s", err.Error())
	}
	
	// 5. 获取每日提币记录 - 每日8点获取
	
	err = c.AddFunc("1 0 8 * * *", func() {
		task.wallet.TransportHis(ctx)
	})
	
	if err != nil {
		task.log.Errorf("配置获取每日提币 cronjob 失败, err: %s", err.Error())
	}
	
	//
	c.Start()
}
