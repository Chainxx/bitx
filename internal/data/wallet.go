package data

import (
	"context"
	"github.com/chainxx/bitx/internal/biz"
)

type walletDataSource struct {
	data *Data
}

func NewWalletDataSource(data *Data) biz.WalletDbRepo {
	return &walletDataSource{
		data: data,
	}
}

/*
data - 记录每日资产信息
*/

func (walletDataSource *walletDataSource) RecordDailyAsset(ctx context.Context, asset biz.DailyAsset) error {
	tx := walletDataSource.data.db.WithContext(ctx).Create(&asset)
	return tx.Error
}

/*
data - 获取每日资产列表信息
*/

func (walletDataSource *walletDataSource) GetDailyAssets(ctx context.Context, time int64) ([]biz.DailyAsset, error) {
	var assets []biz.DailyAsset
	tx := walletDataSource.data.db.WithContext(ctx).
		Where("time = ?", time).
		Find(&assets)
	return assets, tx.Error
}

/*
data - 获取一段时间资产数据
*/

func (walletDataSource *walletDataSource) GetRangeDataByAsset(ctx context.Context, asset string, startTime, endTime int64) ([]biz.DailyAsset, error) {
	var assets []biz.DailyAsset
	
	tx := walletDataSource.data.db.WithContext(ctx).
		Where("asset = ? and time >= ? and time <= ?", asset, startTime, endTime).
		Order("time asc").
		Find(&assets)
	
	return assets, tx.Error
}

/*
data - 记录充值信息
*/

func (walletDataSource *walletDataSource) RecordRecharge(ctx context.Context, recharge biz.Recharge) error {
	tx := walletDataSource.data.db.WithContext(ctx).Create(&recharge)
	return tx.Error
}

/*
data - 获取充值信息
*/

func (walletDataSource *walletDataSource) GetRecharge(ctx context.Context, coin string, startTime, endTime int64) ([]biz.Recharge, error) {
	var recharges []biz.Recharge
	
	tx := walletDataSource.data.db.WithContext(ctx)
	
	if coin != "" {
		tx.Where("coin = ?", coin)
	}
	
	if startTime != 0 {
		tx.Where("insert_time >= ?", startTime)
	}
	
	if endTime != 0 {
		tx.Where("insert_time <= ?", endTime)
	}
	
	tx.Order("insert_time desc")
	tx.Find(&recharges)
	
	return recharges, tx.Error
}

/*
data - 记录传输(提币)记录
*/

func (walletDataSource *walletDataSource) RecordTransport(ctx context.Context, transport biz.Transport) error {
	tx := walletDataSource.data.db.WithContext(ctx).Create(&transport)
	return tx.Error
}
