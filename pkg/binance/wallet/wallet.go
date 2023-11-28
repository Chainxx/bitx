package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startopsz/bitx/pkg/binance/common"
	"github.com/startopsz/bitx/pkg/binance/signed"
	"net/url"
	"time"
)

type SystemStatusInfo struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func SystemStatus() {
	uri := "/sapi/v1/system/status"
	reqUrl := common.Url(uri, "")
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return
	}
	
	var systemStatusInfo SystemStatusInfo
	err = json.Unmarshal(resp, &systemStatusInfo)
	if err != nil {
		return
	}
	
	fmt.Println("systemStatusResp: ", systemStatusInfo)
}

type Coin struct {
	Coin              string        `json:"coin"`
	DepositAllEnable  bool          `json:"depositAllEnable"`
	Free              string        `json:"free"`
	Freeze            string        `json:"freeze"`
	Ipoable           string        `json:"ipoable"`
	Ipoing            string        `json:"ipoing"`
	IsLegalMoney      bool          `json:"isLegalMoney"`
	Locked            string        `json:"locked"`
	Name              string        `json:"name"`
	NetworkList       []CoinNetwork `json:"networkList"`
	Storage           string        `json:"storage"`
	Trading           bool          `json:"trading"`
	WithdrawAllEnable bool          `json:"withdrawAllEnable"`
	Withdrawing       string        `json:"withdrawing"`
}

type CoinNetwork struct {
	AddressRegex            string `json:"addressRegex"`
	Coin                    string `json:"coin"`
	DepositDesc             string `json:"depositDesc"`
	DepositEnable           bool   `json:"depositEnable"`
	IsDefault               bool   `json:"isDefault"`
	MemoRegex               string `json:"memoRegex"`
	MinConfirm              int64  `json:"minConfirm"` // 上账所需的最小确认数
	Name                    string `json:"name"`
	Network                 string `json:"network"`
	ResetAddressStatus      bool   `json:"resetAddressStatus"`
	SpecialTips             string `json:"specialTips"`
	UnLockConfirm           int64  `json:"unLockConfirm"` // 解锁需要的确认数
	WithdrawDesc            string `json:"withdrawDesc"`  // 仅在提现关闭时返回
	WithdrawEnable          bool   `json:"withdrawEnable"`
	WithdrawFee             string `json:"withdrawFee"`
	WithdrawIntegerMultiple string `json:"withdrawIntegerMultiple"`
	WithdrawMax             string `json:"withdrawMax"`
	WithdrawMin             string `json:"withdrawMin"`
	SameAddress             bool   `json:"sameAddress"` // 是否需要memo
	EstimatedArrivalTime    int64  `json:"estimatedArrivalTime"`
	Busy                    bool   `json:"busy"`
}

/*
获取所有的数字货币
*/

func GetAllCoin(apiKey, secretKey string) error {
	uri := "/sapi/v1/capital/config/getall"
	
	values := url.Values{}
	
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	var coins []Coin
	err = json.Unmarshal(resp, &coins)
	if err != nil {
		return err
	}
	
	//fmt.Println("coins: ", coins)
	
	var count int
	for _, coin := range coins {
		
		b, _ := json.Marshal(coin)
		fmt.Println("coin: ", string(b))
		fmt.Println("")
		count += 1
	}
	
	return nil
}

/*
查询每日资产快照 (USER_DATA)

SPOT 现货
MARGIN 杠杆
FUTURES 期货

显示历史拥有过的资产，如果历史拥有，然后清空会显示 0.

每一种类型的 Response Content 不一样

recvWindow: 建议配置在60秒以内，5秒左右
*/

type PerDayAssets struct {
	Code        int64               `json:"code"`
	Msg         string              `json:"msg"`
	SnapshotVos []AssetsSnapshotVos `json:"snapshotVos"`
}

type AssetsSnapshotVos struct {
	Data       AssetsSnapshotVosData `json:"data"`
	Type       string                `json:"type"` // spot  现货
	UpdateTime int64                 `json:"updateTime"`
}

type AssetsSnapshotVosData struct {
	Balances        []AssetsSnapshotVosDataBalance `json:"balances"`
	TotalAssetOfBtc string                         `json:"totalAssetOfBtc"`
}

type AssetsSnapshotVosDataBalance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

func GetPerDayAssets(apiKey, secretKey string, startTime, endTime int64) (PerDayAssets, error) {
	var perDayAssets PerDayAssets
	
	//
	timestamp := time.Now().UnixMilli()
	
	uri := "/sapi/v1/accountSnapshot"
	
	values := url.Values{}
	values.Set("type", "SPOT")
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("startTime", fmt.Sprintf("%d", startTime))
	values.Set("endTime", fmt.Sprintf("%d", endTime))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return perDayAssets, err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	//
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	if err != nil {
		return perDayAssets, err
	}
	
	err = json.Unmarshal(resp, &perDayAssets)
	if err != nil {
		return perDayAssets, err
	}
	
	if perDayAssets.Code != 200 {
		return perDayAssets, errors.New(perDayAssets.Msg)
	}
	
	return perDayAssets, nil
}

/*
 关闭站内划转


*/

func DisableFastWithdrawSwitch() {
	
}

/*
 开启站内划转
*/

func EnableFastWithdrawSwitch() {
	
}

/*
提币
*/

func TransportCoin() {
	
}

/*
获取充值历史

同时提交startTime 与 endTime间隔不得超过90天.
*/

type RechargeHisInfo struct {
	Id            string `json:"id,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Coin          string `json:"coin,omitempty"`
	Network       string `json:"network,omitempty"`
	Status        int    `json:"status,omitempty"`
	Address       string `json:"address,omitempty"`
	AddressTag    string `json:"addressTag,omitempty"`
	TxId          string `json:"txId,omitempty"`
	InsertTime    int64  `json:"insertTime,omitempty"`
	TransferType  int    `json:"transferType,omitempty"`
	ConfirmTimes  string `json:"confirmTimes,omitempty"`
	UnlockConfirm int    `json:"unlockConfirm,omitempty"`
	WalletType    int    `json:"walletType,omitempty"`
}

func RechargeHis(apiKey, secretKey string, startTime, endTime int64) ([]RechargeHisInfo, error) {
	var rechargeHisInfo []RechargeHisInfo
	
	uri := "/sapi/v1/capital/deposit/hisrec"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("startTime", fmt.Sprintf("%d", startTime))
	values.Set("endTime", fmt.Sprintf("%d", endTime))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return rechargeHisInfo, err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return rechargeHisInfo, err
	}
	
	err = json.Unmarshal(resp, &rechargeHisInfo)
	if err != nil {
		return rechargeHisInfo, err
	}
	
	return rechargeHisInfo, nil
}

/*
获取提币历史
*/

type TransportHisInfo struct {
	Id              string `json:"id,omitempty"`              // 该笔提现在币安的id
	Amount          string `json:"amount,omitempty"`          // 提现转出金额
	TransactionFee  string `json:"transactionFee,omitempty"`  // 手续费
	Coin            string `json:"coin,omitempty"`            // 虚拟货币
	Status          int    `json:"status,omitempty"`          // 0:已发送确认Email,1:已被用户取消 2:等待确认 3:被拒绝 4:处理中 5:提现交易失败 6 提现完成
	Address         string `json:"address,omitempty"`         // 地址
	TxId            string `json:"txId,omitempty"`            // 提现交易id
	ApplyTime       string `json:"applyTime,omitempty"`       // UTC 时间
	Network         string `json:"network,omitempty"`         // 网络
	TransferType    int    `json:"transferType,omitempty"`    // 1: 站内转账, 0: 站外转账
	WithdrawOrderId string `json:"withdrawOrderId,omitempty"` // 自定义ID, 如果没有则不返回该字段
	Info            string `json:"info,omitempty"`            // 提币失败原因
	ConfirmNo       int    `json:"confirmNo,omitempty"`       // 提现确认数
	WalletType      int    `json:"walletType,omitempty"`      //1: 资金钱包 0:现货钱包
	TxKey           string `json:"txKey,omitempty"`
	CompleteTime    string `json:"completeTime,omitempty" ` // 提现完成，成功下账时间(UTC)
}

func TransportHis(apiKey, secretKey string, startTime, endTime int64) ([]TransportHisInfo, error) {
	var transportHisInfo []TransportHisInfo
	uri := "/sapi/v1/capital/withdraw/history"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("startTime", fmt.Sprintf("%d", startTime))
	values.Set("endTime", fmt.Sprintf("%d", endTime))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return transportHisInfo, err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return transportHisInfo, err
	}
	
	err = json.Unmarshal(resp, &transportHisInfo)
	
	return transportHisInfo, err
}

/*
 获取充值地址
*/

type CoinAddr struct {
	Address string `json:"address"`
	Coin    string `json:"coin"`
	Tag     string `json:"tag"`
	Url     string `json:"url"`
}

func GetCoinAddr(apiKey, secretKey string, coin string) error {
	uri := "/sapi/v1/capital/deposit/address"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("coin", coin)
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	var coinAddr CoinAddr
	err = json.Unmarshal(resp, &coinAddr)
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", coinAddr)
	return nil
}

/*
账户状态
*/

func AccountStatus(apiKey, secretKey string) error {
	uri := "/sapi/v1/account/status"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", string(resp))
	return nil
}

/*
小额资产转换BNB历史
*/

type ToBNBHisResult struct {
	Total              int                 `json:"total"`
	UserAssetDribblets []UserAssetDribblet `json:"userAssetDribblets"`
}

type UserAssetDribblet struct {
	OperateTime              int64                     `json:"operateTime"`
	TotalServiceChargeAmount string                    `json:"totalServiceChargeAmount"`
	TotalTransferedAmount    string                    `json:"totalTransferedAmount"`
	TransId                  int64                     `json:"transId"`
	UserAssetDribbletDetails []UserAssetDribbletDetail `json:"userAssetDribbletDetails"`
}

type UserAssetDribbletDetail struct {
	FromAsset           string `json:"fromAsset"`
	Amount              string `json:"amount"`
	TransferedAmount    string `json:"transferedAmount"`
	ServiceChargeAmount string `json:"serviceChargeAmount"`
	OperateTime         int64  `json:"operateTime"`
	TransId             int64  `json:"transId"`
}

func ToBNBHis(apiKey, secretKey string) (ToBNBHisResult, error) {
	var toBNBHisResult ToBNBHisResult
	uri := "/sapi/v1/asset/dribblet"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return toBNBHisResult, err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return toBNBHisResult, err
	}
	
	err = json.Unmarshal(resp, &toBNBHisResult)
	return toBNBHisResult, err
}

type ToBNBAsset struct {
	Details []struct {
		Asset            string `json:"asset"`            //资产名
		AssetFullName    string `json:"assetFullName"`    //资产全称
		AmountFree       string `json:"amountFree"`       //可转换数量
		ToBTC            string `json:"toBTC"`            //等值BTC
		ToBNB            string `json:"toBNB"`            //可转换BNB（未扣除手续费）
		ToBNBOffExchange string `json:"toBNBOffExchange"` //可转换BNB（已扣除手续费）
		Exchange         string `json:"exchange"`         //手续费
	} `json:"details"`
	TotalTransferBtc   string `json:"totalTransferBtc"`
	TotalTransferBNB   string `json:"totalTransferBNB"`
	DribbletPercentage string `json:"dribbletPercentage"`
}

/*
获取可以转换成BNB的小额资产
*/

func GetLittleAssetToBNB(apiKey, secretKey string) error {
	uri := "/sapi/v1/asset/dust-btc"
	
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("POST", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	var toBNBAsset ToBNBAsset
	
	err = json.Unmarshal(resp, &toBNBAsset)
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", toBNBAsset)
	return nil
}

type AssetToBNB struct {
	TotalServiceCharge string `json:"totalServiceCharge"`
	TotalTransfered    string `json:"totalTransfered"`
	TransferResult     []struct {
		Amount              string `json:"amount"`
		FromAsset           string `json:"fromAsset"`
		OperateTime         int64  `json:"operateTime"`
		ServiceChargeAmount string `json:"serviceChargeAmount"`
		TranId              int64  `json:"tranId"`
		TransferedAmount    string `json:"transferedAmount"`
	} `json:"transferResult"`
}

/*
小额资产转换
*/

func LittleAssetToBNB(apiKey, secretKey string, asset string) error {
	uri := "/sapi/v1/asset/dust"
	
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("asset", asset)
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("POST", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", string(resp))
	return nil
}

/*
资产利息记录
*/

type AssetDividendInfo struct {
	Id      int64  `json:"id,omitempty"`
	Amount  string `json:"amount,omitempty"`
	Asset   string `json:"asset,omitempty"`
	DivTime int64  `json:"divTime,omitempty"`
	EnInfo  string `json:"enInfo,omitempty"`
	TranId  int64  `json:"tranId,omitempty"`
}

func AssetDividend(apiKey, secretKey string, startTime, endTime int64) error {
	uri := "/sapi/v1/asset/assetDividend"
	
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("startTime", fmt.Sprintf("%d", startTime))
	values.Set("endTime", fmt.Sprintf("%d", endTime))
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", string(resp))
	return nil
}

/*
 上架资产详情
*/

type CoinInfo struct {
	MinWithdrawAmount string `json:"minWithdrawAmount,omitempty"` //最小提现数量
	DepositStatus     bool   `json:"depositStatus,omitempty"`     //是否可以充值(只有所有网络都关闭充值才为false)
	WithdrawFee       string `json:"withdrawFee,omitempty"`       //提现手续费
	WithdrawStatus    bool   `json:"withdrawStatus,omitempty"`    //是否开放提现(只有所有网络都关闭提币才为false)
	DepositTip        string `json:"depositTip,omitempty"`        //暂停充值的原因(如果暂停才有这一项)
}

func AssetDetail(apiKey, secretKey string) error {
	uri := "/sapi/v1/asset/assetDetail"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	coinInfo := make(map[string]CoinInfo)
	
	err = json.Unmarshal(resp, &coinInfo)
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", coinInfo)
	return nil
}

/*
交易手续费率查询
*/

type TradeFee struct {
	Symbol          string `json:"symbol,omitempty"`
	MakerCommission string `json:"makerCommission,omitempty"`
	TakerCommission string `json:"takerCommission,omitempty"`
}

func GetTradeFee(apiKey, secretKey, symbol string) error {
	uri := "/sapi/v1/asset/tradeFee"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("symbol", symbol)
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	var tradeFees []TradeFee
	
	err = json.Unmarshal(resp, &tradeFees)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", tradeFees)
	return nil
}

/*
用户万向划转

MAIN_UMFUTURE 现货钱包转向U本位合约钱包
MAIN_CMFUTURE 现货钱包转向币本位合约钱包
MAIN_MARGIN 现货钱包转向杠杆全仓钱包
UMFUTURE_MAIN U本位合约钱包转向现货钱包
UMFUTURE_MARGIN U本位合约钱包转向杠杆全仓钱包
CMFUTURE_MAIN 币本位合约钱包转向现货钱包
MARGIN_MAIN 杠杆全仓钱包转向现货钱包
MARGIN_UMFUTURE 杠杆全仓钱包转向U本位合约钱包
MARGIN_CMFUTURE 杠杆全仓钱包转向币本位合约钱包
CMFUTURE_MARGIN 币本位合约钱包转向杠杆全仓钱包
ISOLATEDMARGIN_MARGIN 杠杆逐仓钱包转向杠杆全仓钱包
MARGIN_ISOLATEDMARGIN 杠杆全仓钱包转向杠杆逐仓钱包
ISOLATEDMARGIN_ISOLATEDMARGIN 杠杆逐仓钱包转向杠杆逐仓钱包
MAIN_FUNDING 现货钱包转向资金钱包
FUNDING_MAIN 资金钱包转向现货钱包
FUNDING_UMFUTURE 资金钱包转向U本位合约钱包
UMFUTURE_FUNDING U本位合约钱包转向资金钱包
MARGIN_FUNDING 杠杆全仓钱包转向资金钱包
FUNDING_MARGIN 资金钱包转向杠杆全仓钱包
FUNDING_CMFUTURE 资金钱包转向币本位合约钱包
CMFUTURE_FUNDING 币本位合约钱包转向资金钱包

*/

func Transfer() {
	
}

/*
查询用户万向划转历史
*/

func GetTransferHis(apiKey, secretKey string, tType string) error {
	uri := "/sapi/v1/asset/transfer"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("type", tType)
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", string(resp))
	return nil
}

/*
 资金账户
*/

type FundingAsset struct {
	Asset        string `json:"asset"`
	Free         string `json:"free"`
	Locked       string `json:"locked"`
	Freeze       string `json:"freeze"`
	Withdrawing  string `json:"withdrawing"`
	BtcValuation string `json:"btcValuation"`
}

func GetFundingAsset(apiKey, secretKey string) ([]FundingAsset, error) {
	var fundingAssets []FundingAsset
	uri := "/sapi/v1/asset/get-funding-asset"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return fundingAssets, err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("POST", reqUrl, "", headers)
	
	if err != nil {
		return fundingAssets, err
	}
	
	err = json.Unmarshal(resp, &fundingAssets)
	return fundingAssets, err
}

/*
用户持仓
*/

type UserAsset struct {
	Asset        string `json:"asset,omitempty"`
	Free         string `json:"free,omitempty"`
	Locked       string `json:"locked,omitempty"`
	Freeze       string `json:"freeze,omitempty"`
	Withdrawing  string `json:"withdrawing,omitempty"`
	Ipoable      string `json:"ipoable,omitempty"`
	BtcValuation string `json:"btcValuation,omitempty"`
}

func GetUserAsset(apiKey, secretKey string) ([]UserAsset, error) {
	var userAssets []UserAsset
	uri := "/sapi/v3/asset/getUserAsset"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	values.Set("needBtcValuation", "true")
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return userAssets, err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("POST", reqUrl, "", headers)
	
	if err != nil {
		return userAssets, err
	}
	
	err = json.Unmarshal(resp, &userAssets)
	return userAssets, err
}

/*
查询用户API Key权限
*/

type APIKeyPerm struct {
	IpRestrict                     bool  `json:"ipRestrict,omitempty"`
	CreateTime                     int64 `json:"createTime,omitempty"`
	EnableWithdrawals              bool  `json:"enableWithdrawals,omitempty"`
	EnableInternalTransfer         bool  `json:"enableInternalTransfer,omitempty"`
	PermitsUniversalTransfer       bool  `json:"permitsUniversalTransfer,omitempty"`
	EnableVanillaOptions           bool  `json:"enableVanillaOptions,omitempty"`
	EnableReading                  bool  `json:"enableReading,omitempty"`
	EnableFutures                  bool  `json:"enableFutures,omitempty"`
	EnableMargin                   bool  `json:"enableMargin,omitempty"`
	EnableSpotAndMarginTrading     bool  `json:"enableSpotAndMarginTrading,omitempty"`
	TradingAuthorityExpirationTime int64 `json:"tradingAuthorityExpirationTime,omitempty"`
}

/*
查询用户API Key权限
*/

func GetAPIKeyPerm(apiKey, secretKey string) error {
	uri := "/sapi/v1/account/apiRestrictions"
	values := url.Values{}
	timestamp := time.Now().UnixMilli()
	
	values.Set("timestamp", fmt.Sprint(timestamp))
	
	signature, err := signed.Sign(secretKey, values.Encode())
	if err != nil {
		return err
	}
	
	values.Set("signature", signature)
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	
	if err != nil {
		return err
	}
	
	var apiKeyPerm APIKeyPerm
	err = json.Unmarshal(resp, &apiKeyPerm)
	if err != nil {
		return err
	}
	
	fmt.Println("resp: ", apiKeyPerm)
	return nil
}
