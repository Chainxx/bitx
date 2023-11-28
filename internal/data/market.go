package data

import (
	"context"
	"github.com/chainxx/bitx/internal/biz"
)

type marketDataResource struct {
	data *Data
}

func NewMarketDataResource(data *Data) biz.MarketDbRepo {
	return &marketDataResource{
		data: data,
	}
}

/*
data - 记录价格
*/

func (marketDataResource *marketDataResource) RecordTickerPrice(ctx context.Context, tickerPrice biz.TickerPrice) error {
	tx := marketDataResource.data.db.WithContext(ctx).Create(&tickerPrice)
	return tx.Error
}

/*
data - 获取Symbol价格
*/

func (marketDataResource *marketDataResource) GetTickerPriceBySymbol(ctx context.Context, symbol string, startTime int64, endTime int64) ([]biz.TickerPrice, error) {
	var tickerPrices []biz.TickerPrice
	tx := marketDataResource.data.db.WithContext(ctx).
		Where("symbol = ? and open_time >= ? and close_time <= ?", symbol, startTime, endTime).
		Order("open_time desc").
		Find(&tickerPrices)
	
	return tickerPrices, tx.Error
}

/*
data - 记录每日价格
*/

func (marketDataResource *marketDataResource) RecordTickerPriceDaily(ctx context.Context, tickerPriceDaily biz.TickerPriceDaily) error {
	tx := marketDataResource.data.db.WithContext(ctx).Create(&tickerPriceDaily)
	return tx.Error
}
