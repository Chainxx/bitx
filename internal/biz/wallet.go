package biz

import (
	"context"
	"fmt"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/chainxx/bitx/pkg/binance/wallet"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

// 每日资产

type DailyAsset struct {
	Asset      string  `json:"asset,omitempty"`
	Free       float64 `json:"free,omitempty"`
	Locked     float64 `json:"locked,omitempty"`
	Time       int64   `json:"time,omitempty"`
	CreateTime int64   `json:"createTime,omitempty"`
}

func (dailyAsset DailyAsset) TableName() string {
	return "daily_asset"
}

// 充值

type Recharge struct {
	Id            string  `json:"id,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	Coin          string  `json:"coin,omitempty"`
	Network       string  `json:"network,omitempty"`
	Status        int     `json:"status,omitempty"`
	Address       string  `json:"address,omitempty"`
	AddressTag    string  `json:"addressTag,omitempty"`
	TxId          string  `json:"txId,omitempty"`
	InsertTime    int64   `json:"insertTime,omitempty"`
	TransferType  int     `json:"transferType,omitempty"`
	ConfirmTimes  string  `json:"confirmTimes,omitempty"`
	UnlockConfirm int     `json:"unlockConfirm,omitempty"`
	WalletType    int     `json:"walletType,omitempty"`
}

func (recharge Recharge) TableName() string {
	return "recharge"
}

// 提币记录

type Transport struct {
	Id              string  `json:"id,omitempty"`
	Amount          float64 `json:"amount,omitempty"`
	TransactionFee  float64 `json:"transactionFee,omitempty"`
	Coin            string  `json:"coin,omitempty"`
	Status          int     `json:"status,omitempty"`
	Address         string  `json:"address,omitempty"`
	TxId            string  `json:"txId,omitempty"`
	ApplyTime       string  `json:"applyTime,omitempty"`
	Network         string  `json:"network,omitempty"`
	TransferType    int     `json:"transferType,omitempty"`
	WithdrawOrderId string  `json:"withdrawOrderId,omitempty"`
	Info            string  `json:"info,omitempty"`
	ConfirmNo       int     `json:"confirmNo,omitempty"`
	WalletType      int     `json:"walletType,omitempty"`
	TxKey           string  `json:"txKey,omitempty"`
	CompleteTime    string  `json:"completeTime,omitempty" `
}

func (transport Transport) TableName() string {
	return "transport"
}

type WalletDbRepo interface {
	RecordDailyAsset(ctx context.Context, asset DailyAsset) error
	GetDailyAssets(ctx context.Context, time int64) ([]DailyAsset, error)
	GetRangeDataByAsset(ctx context.Context, asset string, startTime, endTime int64) ([]DailyAsset, error)
	
	RecordRecharge(ctx context.Context, recharge Recharge) error
	GetRecharge(ctx context.Context, coin string, startTime, endTime int64) ([]Recharge, error)
	
	RecordTransport(ctx context.Context, transport Transport) error
}

type WalletUseCase struct {
	ApiKey       string
	SecretKey    string
	walletDbRepo WalletDbRepo
	log          *log.Helper
}

func NewWalletUseCase(binance *conf.Binance, walletDbRepo WalletDbRepo, logger log.Logger) *WalletUseCase {
	return &WalletUseCase{
		walletDbRepo: walletDbRepo,
		ApiKey:       binance.ApiKey,
		SecretKey:    binance.SecretKey,
		log:          log.NewHelper(logger),
	}
}

/*
biz - 获取系统状态 - 用于检查
*/

func (walletUseCase *WalletUseCase) SystemStatus() {

}

func (walletUseCase *WalletUseCase) GetPerDayAssetsByDB(ctx context.Context, asset string, startTime int64, endTime int64) ([]DailyAsset, error) {
	return walletUseCase.walletDbRepo.GetRangeDataByAsset(ctx, asset, startTime, endTime)
}

/*
获取每日资产

可以每日报告资产信息 - UTC 时间
*/

func (walletUseCase *WalletUseCase) GetPerDayAssets(ctx context.Context) {
	endTime := time.Now().Unix() * 1000
	startTime := time.Now().Add(-1*24*time.Hour).Unix() * 1000
	walletUseCase.log.Infof("开始执行每日资产获取数据任务, startTime: %d, endTime: %d", startTime, endTime)
	
	perDayAssets, err := wallet.GetPerDayAssets(walletUseCase.ApiKey, walletUseCase.SecretKey, startTime, endTime)
	if err != nil {
		walletUseCase.log.Errorf("获取每日资产失败, err: %s", err.Error())
		return
	}
	
	if perDayAssets.Code != 200 {
		walletUseCase.log.Errorf("获取每日资产失败, err: %s", perDayAssets.Msg)
		return
	}
	
	if len(perDayAssets.SnapshotVos) == 0 {
		walletUseCase.log.Infof("获取每日资产, 资产列表为空, 无资产, startTime: %d, endTime: %d", startTime, endTime)
		return
	}
	
	for _, vos := range perDayAssets.SnapshotVos {
		//fmt.Printf("Time: %s, TotalAssetOfBtc: %s.\n", timestamp.Time(vos.UpdateTime).String(), vos.Data.TotalAssetOfBtc)
		for _, balance := range vos.Data.Balances {
			
			free, err := strconv.ParseFloat(balance.Free, 32)
			if err != nil {
				walletUseCase.log.Errorf("记录每日资产信息失败, 将 balance free (%s) 转换 float32 失败, err: ", balance.Free, err.Error())
				continue
			}
			
			locked, err := strconv.ParseFloat(balance.Locked, 32)
			if err != nil {
				walletUseCase.log.Errorf("记录每日资产信息失败, 将 balance locked (%s) 转换 float32 失败, err: ", balance.Locked, err.Error())
				continue
			}
			
			err = walletUseCase.walletDbRepo.RecordDailyAsset(ctx, DailyAsset{
				Asset:      balance.Asset,
				Free:       free,
				Locked:     locked,
				Time:       vos.UpdateTime / 1000,
				CreateTime: time.Now().Unix(),
			})
			
			if err != nil {
				walletUseCase.log.Errorf("记录每日资产信息失败, Asset: %s, Free: %s, Locked: %s",
					balance.Asset, balance.Free, balance.Locked)
			}
		}
	}
}

/*
获取充值历史
*/

func (walletUseCase *WalletUseCase) RechargeHis(ctx context.Context) {
	endTime := time.Now().Unix() * 1000
	startTime := time.Now().Add(-1*24*time.Hour).Unix() * 1000
	
	walletUseCase.log.Infof("开始获取历史充值信息, startTime: %d, endTime: %d", startTime, endTime)
	
	rechargeHisInfos, err := wallet.RechargeHis(walletUseCase.ApiKey, walletUseCase.SecretKey, startTime, endTime)
	if err != nil {
		walletUseCase.log.Errorf("获取充值历史数据失败, err: %s", err.Error())
		return
	}
	
	if len(rechargeHisInfos) == 0 {
		walletUseCase.log.Infof("获取充值历史数据, 历史充值为空, 无充值信息, startTime: %d, endTime: %d", startTime, endTime)
		return
	}
	
	for _, rechargeHis := range rechargeHisInfos {
		
		amount, err := strconv.ParseFloat(rechargeHis.Amount, 32)
		if err != nil {
			walletUseCase.log.Errorf("获取历史充值信息，转换充值信息amount(%s)为float64失败, id: %s, err: %s",
				rechargeHis.Amount, rechargeHis.Id, err.Error())
			continue
		}
		
		err = walletUseCase.walletDbRepo.RecordRecharge(ctx, Recharge{
			Id:            rechargeHis.Id,
			Amount:        amount,
			Coin:          rechargeHis.Coin,
			Network:       rechargeHis.Network,
			Status:        rechargeHis.Status,
			Address:       rechargeHis.Address,
			AddressTag:    rechargeHis.AddressTag,
			TxId:          rechargeHis.TxId,
			InsertTime:    rechargeHis.InsertTime / 1000,
			TransferType:  rechargeHis.TransferType,
			ConfirmTimes:  rechargeHis.ConfirmTimes,
			UnlockConfirm: rechargeHis.UnlockConfirm,
		})
		
		if err != nil {
			walletUseCase.log.Errorf("记录历史充值信息失败, id: %s, err: %s", rechargeHis.Id, err.Error())
		}
		
	}
}

/*
获取提币历史
*/

func (walletUseCase *WalletUseCase) TransportHis(ctx context.Context) {
	endTime := time.Now().Unix() * 1000
	startTime := time.Now().Add(-1*24*time.Hour).Unix() * 1000
	
	walletUseCase.log.Infof("开始获取提币历史, startTime: %d, endTime: %d.", startTime, endTime)
	
	transportHis, err := wallet.TransportHis(walletUseCase.ApiKey, walletUseCase.SecretKey, startTime, endTime)
	if err != nil {
		walletUseCase.log.Errorf("获取提币历史失败, err: %s", err.Error())
		return
	}
	
	if len(transportHis) == 0 {
		walletUseCase.log.Errorf("获取提币历史, 提币历史为空, 无提币记录")
		return
	}
	
	for _, transport := range transportHis {
		
		amount, err := strconv.ParseFloat(transport.Amount, 32)
		if err != nil {
			walletUseCase.log.Errorf("获取提币历史，转换充值信息amount(%s)为float64失败, id: %s, err: %s",
				transport.Amount, transport.Id, err.Error())
			continue
		}
		
		transactionFee, err := strconv.ParseFloat(transport.TransactionFee, 32)
		if err != nil {
			walletUseCase.log.Errorf("获取提币历史，转换充值信息TransactionFee(%s)为float64失败, id: %s, err: %s",
				transport.TransactionFee, transport.Id, err.Error())
			continue
		}
		
		err = walletUseCase.walletDbRepo.RecordTransport(ctx, Transport{
			Id:              transport.Id,
			Amount:          amount,
			TransactionFee:  transactionFee,
			Coin:            transport.Coin,
			Status:          transport.Status,
			Address:         transport.Address,
			TxId:            transport.TxId,
			ApplyTime:       transport.ApplyTime,
			Network:         transport.Network,
			TransferType:    transport.TransferType,
			WithdrawOrderId: transport.WithdrawOrderId,
			Info:            transport.Info,
			ConfirmNo:       transport.ConfirmNo,
			WalletType:      transport.WalletType,
			TxKey:           transport.TxKey,
			CompleteTime:    transport.CompleteTime,
		})
		
		if err != nil {
			walletUseCase.log.Errorf("获取提币历史, 记录数据到数据库失败, id: %s, txId: %s, err: %s",
				transport.Id, transport.TxId, err.Error())
		}
	}
}

/*
biz - 账户状态
*/

func (walletUseCase *WalletUseCase) AccountStatus(ctx context.Context) {
	err := wallet.AccountStatus(walletUseCase.ApiKey, walletUseCase.SecretKey)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
}

/*
小额资产转换BNB历史
*/

func (walletUseCase *WalletUseCase) ToBNBHis(ctx context.Context) {
	toBNBHisResult, err := wallet.ToBNBHis(walletUseCase.ApiKey, walletUseCase.SecretKey)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	fmt.Println("toBNBHisResult: ", toBNBHisResult)
}

/*
资金账户 - 在使用杠杆是，需要资金账户有对应的资产，如果不使用杠杆等功能，不需要关注资金账户。
*/

func (walletUseCase *WalletUseCase) GetFundingAsset(ctx context.Context) {
	fundingAssets, err := wallet.GetFundingAsset(walletUseCase.ApiKey, walletUseCase.SecretKey)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	for _, fundingAsset := range fundingAssets {
		fmt.Println("fundingAsset: ", fundingAsset)
	}
}

/*
用户持仓
*/

func (walletUseCase *WalletUseCase) GetUserAsset(ctx context.Context) {
	userAssets, err := wallet.GetUserAsset(walletUseCase.ApiKey, walletUseCase.SecretKey)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	for _, userAsset := range userAssets {
		fmt.Println("userAsset: ", userAsset)
	}
}
