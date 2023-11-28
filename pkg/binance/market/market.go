package market

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startopsz/bitx/pkg/binance/common"
	"net/url"
)

/*
测试服务器连通性
*/

func Ping() error {
	uri := "/api/v3/ping"
	reqUrl := common.Url(uri, "")
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(resp))
	
	return nil
}

/*
获取服务器时间
*/

func ServerTime() error {
	uri := "/api/v3/time"
	reqUrl := common.Url(uri, "")
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(resp))
	
	return nil
}

/*
交易规范信息
*/

func ExchangeInfo() {
	
}

/*
深度信息
*/

type DepthInfo struct {
	LastUpdateId int64      `json:"lastUpdateId,omitempty"`
	Bids         [][]string `json:"bids,omitempty"` // 价位 & 挂单量
	Asks         [][]string `json:"asks,omitempty"`
}

func Depth() error {
	uri := "/api/v3/depth"
	
	values := url.Values{}
	values.Add("symbol", "BNBBTC")
	
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return err
	}
	
	var depthInfo DepthInfo
	
	err = json.Unmarshal(resp, &depthInfo)
	if err != nil {
		return err
	}
	
	fmt.Println("depthInfo: ", depthInfo)
	
	return nil
}

/*
近期成交列表
*/

type TradeInfo struct {
	Id           int64  `json:"id,omitempty"`
	Price        string `json:"price,omitempty"`        // 价格
	Qty          string `json:"qty,omitempty"`          // 成交数量
	QuoteQty     string `json:"quoteQty,omitempty"`     // 报价数量
	Time         int64  `json:"time,omitempty"`         // 交易成交时间, 和websocket中的T一致.
	IsBuyerMaker bool   `json:"isBuyerMaker,omitempty"` // 是否购买者市场
	IsBestMatch  bool   `json:"isBestMatch,omitempty"`  // 是否最佳匹配
}

func Trades(symbol string) ([]TradeInfo, error) {
	var tradesInfo []TradeInfo
	uri := "/api/v3/trades"
	
	values := url.Values{}
	values.Add("symbol", symbol)
	//values.Add("symbol", "BNBUSDT")
	
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return tradesInfo, err
	}
	
	err = json.Unmarshal(resp, &tradesInfo)
	if err != nil {
		return tradesInfo, err
	}
	
	return tradesInfo, nil
}

/*
查询历史成交
*/

func TradesHis(apiKey string, symbol string) ([]TradeInfo, error) {
	var tradeHis []TradeInfo
	uri := "/api/v3/historicalTrades"
	
	values := url.Values{}
	values.Add("symbol", symbol)
	//values.Add("symbol", "BNBUSDT")
	
	reqUrl := common.Url(uri, values.Encode())
	
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = apiKey
	
	resp, _, err := common.Do("GET", reqUrl, "", headers)
	if err != nil {
		return tradeHis, err
	}
	
	err = json.Unmarshal(resp, &tradeHis)
	if err != nil {
		return tradeHis, err
	}
	
	return tradeHis, nil
}

/*
K 线
*/

func Kline() error {
	uri := "/api/v3/klines"
	
	values := url.Values{}
	values.Set("symbol", "BNBUSDT")
	values.Set("interval", "")
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(resp))
	
	return nil
}

/*
当前平均价格
*/

type AvgPriceResult struct {
	Mins  int    `json:"mins"`  // 分钟
	Price string `json:"price"` // 价格
}

func AvgPrice(symbol string) (AvgPriceResult, error) {
	var avgPriceResult AvgPriceResult
	uri := "/api/v3/avgPrice"
	
	values := url.Values{}
	values.Set("symbol", symbol)
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return avgPriceResult, err
	}
	
	err = json.Unmarshal(resp, &avgPriceResult)
	
	return avgPriceResult, err
	
}

// (prevClosePrice 参数用于指示在某一特定交易对（如BTC/USDT）上的前一个交易周期（通常是前一个交易日、前一个小时或前一个分钟，取决于所选择的K线间隔）的收盘价格。)

type Ticker struct {
	Symbol             string `json:"symbol,omitempty"`             // 币种
	PriceChange        string `json:"priceChange,omitempty"`        // 价格变化
	PriceChangePercent string `json:"priceChangePercent,omitempty"` // 价格变化百分比
	WeightedAvgPrice   string `json:"weightedAvgPrice,omitempty"`   // 加权平均价格
	PrevClosePrice     string `json:"prevClosePrice,omitempty"`     // 前一交易周期的收盘价格
	BidPrice           string `json:"bidPrice,omitempty"`           // 出价价格
	BidQty             string `json:"bidQty,omitempty"`             // 出价数量
	AskPrice           string `json:"askPrice,omitempty"`           // 询价价格
	AskQty             string `json:"askQty,omitempty"`             // 询价数量
	OpenPrice          string `json:"openPrice,omitempty"`          // 开始价格
	HighPrice          string `json:"highPrice,omitempty"`          // 高价
	LowPrice           string `json:"lowPrice,omitempty"`           // 低价
	LastPrice          string `json:"lastPrice,omitempty"`          // 最后价格
	LastQty            string `json:"lastQty,omitempty"`            // 最后价格数量
	Volume             string `json:"volume,omitempty"`             // 体量
	QuoteVolume        string `json:"quoteVolume,omitempty"`        // 报价单数量
	OpenTime           int64  `json:"openTime,omitempty"`           // 开放时间
	CloseTime          int64  `json:"closeTime,omitempty"`          // 关闭时间
	FirstId            int64  `json:"firstId,omitempty"`            //
	LastId             int64  `json:"lastId,omitempty"`             //
	Count              int64  `json:"count,omitempty"`              // 统计时间内交易笔数
}

/*
24hr 价格变动情况
非完整的一天 00:00 ~ 23:59:59, 请求时间: CloseTime, 24h 开始时间为请求时间往前推算24小时
*/

func DayTicker(symbol string) (Ticker, error) {
	var ticker Ticker
	uri := "/api/v3/ticker/24hr"
	
	values := url.Values{}
	values.Set("symbol", symbol)
	//values.Set("symbol", "BNBUSDT")
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return ticker, err
	}
	
	err = json.Unmarshal(resp, &ticker)
	if err != nil {
		return ticker, err
	}
	
	return ticker, nil
}

/*
最新价格
*/

type Price struct {
	Symbol string `json:"symbol,omitempty"`
	Price  string `json:"price,omitempty"`
}

func LastPrice(symbol string) (Price, error) {
	var price Price
	uri := "/api/v3/ticker/price"
	
	values := url.Values{}
	values.Add("symbol", symbol)
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return price, err
	}
	
	err = json.Unmarshal(resp, &price)
	return price, err
}

/*
当前最优挂单
*/

type GoodTickerInfo struct {
	Symbol   string `json:"symbol,omitempty"`   // 币种
	BidPrice string `json:"bidPrice,omitempty"` // 买方愿意出的最高价
	BidQty   string `json:"bidQty,omitempty"`   // 买方买的数量
	AskPrice string `json:"askPrice,omitempty"` // 卖方接受的价格
	AskQty   string `json:"askQty,omitempty"`   // 卖方卖的数量
}

func GoodTicker(symbol string) (GoodTickerInfo, error) {
	var goodTickerInfo GoodTickerInfo
	uri := "/api/v3/ticker/bookTicker"
	
	values := url.Values{}
	values.Add("symbol", symbol)
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return goodTickerInfo, err
	}
	
	err = json.Unmarshal(resp, &goodTickerInfo)
	if err != nil {
		return goodTickerInfo, err
	}
	
	return goodTickerInfo, nil
}

/*
滚动窗口价格变动统计

重要参考

windowSize: type (ENUM)

默认为 1d
windowSize 支持的值:
如果是分钟: 1m,2m....59m
如果是小时: 1h, 2h....23h
如果是天: 1d...7d

不可以组合使用, 比如1d2h

个人理解为: 一段时间内(该时间需要符合上述枚举需要)价格变动

Important: 成交量 和 成交笔数并不是一个含义
*/

type TickerPriceInfo struct {
	Symbol             string `json:"symbol,omitempty"`
	PriceChange        string `json:"priceChange,omitempty"`        // 价格变化
	PriceChangePercent string `json:"priceChangePercent,omitempty"` // 价格变化百分比
	WeightedAvgPrice   string `json:"weightedAvgPrice,omitempty"`   // 加权平均价格
	OpenPrice          string `json:"openPrice,omitempty"`          // 开盘价, 开始价格
	HighPrice          string `json:"highPrice,omitempty"`          // 高价
	LowPrice           string `json:"lowPrice,omitempty"`           // 低价
	LastPrice          string `json:"lastPrice,omitempty"`          // 最后价格
	Volume             string `json:"volume,omitempty"`             // 成交量 (BTC/USDT: 多少个BTC被买入和卖出)
	QuoteVolume        string `json:"quoteVolume,omitempty"`        // 报价货币成交量 (多少个USDT的价值被交易)
	OpenTime           int64  `json:"openTime,omitempty"`           // ticker的开始时间
	CloseTime          int64  `json:"closeTime,omitempty"`          // ticker的结束时间
	FirstId            int64  `json:"firstId,omitempty"`            // 统计时间内的第一笔trade id
	LastId             int64  `json:"lastId,omitempty"`             //
	Count              int64  `json:"count,omitempty"`              // 统计时间内交易笔数
}

/*
windowsSize 秒
计算结果为除法向下取整，也就是7天23小时认为是7天， 1天23小时认为是1天

windowsSize 默认为 1d
windowSize 支持的值:
如果是分钟: 1m,2m....59m
如果是小时: 1h, 2h....23h
如果是天: 1d...7d
*/

type WindowsSize int64

func (windowsSize WindowsSize) String() (string, error) {
	minute := windowsSize / 60
	if minute == 0 {
		return "", errors.New(fmt.Sprintf("windowsSize: %d, too low", windowsSize))
	}
	
	if minute < 59 {
		return fmt.Sprintf("%dm", minute), nil
	}
	
	hour := windowsSize / 60 / 60
	if hour < 23 {
		return fmt.Sprintf("%dh", hour), nil
	}
	
	day := windowsSize / 60 / 60 / 24
	if day <= 7 {
		return fmt.Sprintf("%dd", day), nil
	}
	//
	return "", errors.New(fmt.Sprintf("windowsSize: %d, too large", windowsSize))
}

func TickerPrice(symbol string, windowsSize WindowsSize) (TickerPriceInfo, error) {
	var tickerPriceInfo TickerPriceInfo
	uri := "/api/v3/ticker"
	
	values := url.Values{}
	values.Set("symbol", symbol)
	
	ws, err := windowsSize.String()
	if err != nil {
		return tickerPriceInfo, err
	}
	
	values.Set("windowSize", ws)
	//values.Set("windowSize", "1m")
	reqUrl := common.Url(uri, values.Encode())
	
	resp, _, err := common.Do("GET", reqUrl, "", nil)
	if err != nil {
		return tickerPriceInfo, err
	}
	
	err = json.Unmarshal(resp, &tickerPriceInfo)
	return tickerPriceInfo, err
}
