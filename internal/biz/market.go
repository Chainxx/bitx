package biz

import (
	"context"
	"fmt"
	"github.com/chainxx/bitx/internal/conf"
	"github.com/chainxx/bitx/pkg/binance/market"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/startopsz/rule/pkg/timestamp"
	"strconv"
	"time"
)

type TickerPrice struct {
	Symbol             string  `json:"symbol,omitempty"`             // 代号
	PriceChange        float64 `json:"priceChange,omitempty"`        // 价格变化
	PriceChangePercent float64 `json:"priceChangePercent,omitempty"` // 价格变化百分比
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,omitempty"`   // 加权平均价格
	OpenPrice          float64 `json:"openPrice,omitempty"`          // 开盘价, 开始价格
	HighPrice          float64 `json:"highPrice,omitempty"`          // 高价
	LowPrice           float64 `json:"lowPrice,omitempty"`           // 低价
	LastPrice          float64 `json:"lastPrice,omitempty"`          // 最后价格
	Volume             float64 `json:"volume,omitempty"`             // 成交量 (BTC/USDT: 多少个BTC被买入和卖出)
	QuoteVolume        float64 `json:"quoteVolume,omitempty"`        // 报价货币成交量 (多少个USDT的价值被交易)
	OpenTime           int64   `json:"openTime,omitempty"`           // ticker的开始时间
	CloseTime          int64   `json:"closeTime,omitempty"`          // ticker的结束时间
	FirstId            int64   `json:"firstId,omitempty"`            // 统计时间内的第一笔trade id
	LastId             int64   `json:"lastId,omitempty"`             //
	Count              int64   `json:"count,omitempty"`              // 统计时间内交易笔数
}

func (tickerPrice TickerPrice) TableName() string {
	return "ticker_price"
}

type TickerPriceDaily struct {
	Symbol             string  `json:"symbol,omitempty"`             // 代号
	PriceChange        float64 `json:"priceChange,omitempty"`        // 价格变化
	PriceChangePercent float64 `json:"priceChangePercent,omitempty"` // 价格变化百分比
	WeightedAvgPrice   float64 `json:"weightedAvgPrice,omitempty"`   // 加权平均价格
	PrevClosePrice     float64 `json:"prevClosePrice,omitempty"`     //
	BidPrice           float64 `json:"bidPrice,omitempty"`           // 出价价格
	BidQty             float64 `json:"bidQty,omitempty"`             // 出价数量
	AskPrice           float64 `json:"askPrice,omitempty"`           // 询价价格
	AskQty             float64 `json:"askQty,omitempty"`             // 询价数量
	OpenPrice          float64 `json:"openPrice,omitempty"`          // 开盘价, 开始价格
	HighPrice          float64 `json:"highPrice,omitempty"`          // 高价
	LowPrice           float64 `json:"lowPrice,omitempty"`           // 低价
	LastPrice          float64 `json:"lastPrice,omitempty"`          // 最后价格
	LastQty            float64 `json:"lastQty,omitempty"`            // 最后价格数量
	Volume             float64 `json:"volume,omitempty"`             // 成交量 (BTC/USDT: 多少个BTC被买入和卖出)
	QuoteVolume        float64 `json:"quoteVolume,omitempty"`        // 报价货币成交量 (多少个USDT的价值被交易)
	OpenTime           int64   `json:"openTime,omitempty"`           // ticker的开始时间
	CloseTime          int64   `json:"closeTime,omitempty"`          // ticker的结束时间
	FirstId            int64   `json:"firstId,omitempty"`            // 统计时间内的第一笔trade id
	LastId             int64   `json:"lastId,omitempty"`             //
	Count              int64   `json:"count,omitempty"`              // 统计时间内交易笔数
}

func (tickerPriceDaily TickerPriceDaily) TableName() string {
	return "ticker_price_daily"
}

type MarketDbRepo interface {
	RecordTickerPrice(ctx context.Context, tickerPrice TickerPrice) error
	GetTickerPriceBySymbol(ctx context.Context, symbol string, startTime, endTime int64) ([]TickerPrice, error)
	RecordTickerPriceDaily(ctx context.Context, tickerPriceDaily TickerPriceDaily) error
}

type TradesUseCase struct {
	ApiKey       string `json:"apiKey"`
	marketDbRepo MarketDbRepo
	symbols      []string
	log          *log.Helper
}

func NewTradesUseCase(binance *conf.Binance, marketDbRepo MarketDbRepo, logger log.Logger) *TradesUseCase {
	return &TradesUseCase{
		marketDbRepo: marketDbRepo,
		ApiKey:       binance.ApiKey,
		symbols:      binance.Symbols,
		log:          log.NewHelper(logger),
	}
}

func (tradesUseCase *TradesUseCase) ListSymbols(ctx context.Context) []string {
	return tradesUseCase.symbols
}

/*
biz - 近期成交列表, 某种代币兑换符号的近期成交信息
*/

func (tradesUseCase *TradesUseCase) Trades(ctx context.Context) {
	TradeInfo, err := market.Trades("BNBUSDT")
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
	
	for _, trade := range TradeInfo {
		fmt.Printf("Time: %s, 成交数量: %s, 报价数量: %s, 价格: %s, 是否购买者市场: %t, 是否最佳匹配: %t.\n",
			timestamp.Time(trade.Time).In(time.Local).String(), trade.Qty, trade.QuoteQty, trade.Price, trade.IsBuyerMaker, trade.IsBestMatch)
	}
}

/*
biz - 查询历史成交
*/

func (tradesUseCase *TradesUseCase) TradesHis(ctx context.Context) {
	tradeHis, err := market.TradesHis(tradesUseCase.ApiKey, "BNBUSDT")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	for _, trade := range tradeHis {
		fmt.Printf("time: %s,  成交数量: %s, 报价数量: %s, 价格: %s, 是否购买者市场: %t, 是否最佳匹配: %t.\n",
			timestamp.Time(trade.Time).In(time.Local).String(), trade.Qty, trade.QuoteQty, trade.Price, trade.IsBuyerMaker, trade.IsBestMatch)
	}
}

/*
biz - 24hr 价格变动情况
非完整的一天 00:00 ~ 23:59:59, 请求时间: CloseTime, 24h 开始时间为请求时间往前推算24小时
*/

func (tradesUseCase *TradesUseCase) DayTicker(ctx context.Context, symbol string) {
	ticker, err := market.DayTicker(symbol)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动接口失败, err: ", symbol, err.Error())
		return
	}
	
	priceChangePercent, err := strconv.ParseFloat(ticker.PriceChangePercent, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换更改百分比失败(PriceChangePercent: %s), err: %s",
			symbol, ticker.PriceChangePercent, err.Error())
		return
	}
	
	weightedAvgPrice, err := strconv.ParseFloat(ticker.WeightedAvgPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(WeightedAvgPrice: %s), err: %s",
			symbol, ticker.WeightedAvgPrice, err.Error())
		return
	}
	
	openPrice, err := strconv.ParseFloat(ticker.OpenPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(OpenPrice: %s), err: %s",
			symbol, ticker.OpenPrice, err.Error())
		return
	}
	
	highPrice, err := strconv.ParseFloat(ticker.HighPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(HighPrice: %s), err: %s",
			symbol, ticker.WeightedAvgPrice, err.Error())
		return
	}
	
	lowPrice, err := strconv.ParseFloat(ticker.LowPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(LowPrice: %s), err: %s",
			symbol, ticker.LowPrice, err.Error())
		return
	}
	
	lastPrice, err := strconv.ParseFloat(ticker.LastPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(LastPrice: %s), err: %s",
			symbol, ticker.LastPrice, err.Error())
		return
	}
	
	volume, err := strconv.ParseFloat(ticker.Volume, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(Volume: %s), err: %s",
			symbol, ticker.Volume, err.Error())
		return
	}
	
	quoteVolume, err := strconv.ParseFloat(ticker.QuoteVolume, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(QuoteVolume: %s), err: %s",
			symbol, ticker.QuoteVolume, err.Error())
		return
	}
	
	priceChange, err := strconv.ParseFloat(ticker.PriceChange, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(PriceChange: %s), err: %s",
			symbol, ticker.PriceChange, err.Error())
		return
	}
	
	prevClosePrice, err := strconv.ParseFloat(ticker.PrevClosePrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(PrevClosePrice: %s), err: %s",
			symbol, ticker.PrevClosePrice, err.Error())
		return
	}
	
	bidPrice, err := strconv.ParseFloat(ticker.BidPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(BidPrice: %s), err: %s",
			symbol, ticker.BidPrice, err.Error())
		return
	}
	
	bidQty, err := strconv.ParseFloat(ticker.BidQty, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(BidQty: %s), err: %s",
			symbol, ticker.BidQty, err.Error())
		return
	}
	
	askPrice, err := strconv.ParseFloat(ticker.AskPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(AskPrice: %s), err: %s",
			symbol, ticker.AskPrice, err.Error())
		return
	}
	
	askQty, err := strconv.ParseFloat(ticker.AskQty, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(AskQty: %s), err: %s",
			symbol, ticker.AskQty, err.Error())
		return
	}
	
	lastQty, err := strconv.ParseFloat(ticker.LastQty, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 转换失败(LastQty: %s), err: %s",
			symbol, ticker.LastQty, err.Error())
		return
	}
	
	if openPrice == 0 || highPrice == 0 {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 返回内容检查失败, openPrice: %d, openTime: %d",
			ticker.Symbol, openPrice, ticker.OpenTime/1000)
		return
	}
	
	err = tradesUseCase.marketDbRepo.RecordTickerPriceDaily(ctx, TickerPriceDaily{
		Symbol:             ticker.Symbol,
		PriceChange:        priceChange,
		PriceChangePercent: priceChangePercent,
		WeightedAvgPrice:   weightedAvgPrice,
		PrevClosePrice:     prevClosePrice,
		BidPrice:           bidPrice,
		BidQty:             bidQty,
		AskPrice:           askPrice,
		AskQty:             askQty,
		OpenPrice:          openPrice,
		HighPrice:          highPrice,
		LowPrice:           lowPrice,
		LastPrice:          lastPrice,
		LastQty:            lastQty,
		Volume:             volume,
		QuoteVolume:        quoteVolume,
		OpenTime:           ticker.OpenTime,
		CloseTime:          ticker.CloseTime,
		FirstId:            ticker.FirstId,
		LastId:             ticker.LastId,
		Count:              ticker.Count,
	})
	
	if err != nil {
		tradesUseCase.log.Errorf("获取(symbol: %s)24小时价格变动, 记录数据到数据库失败, err: %s",
			symbol, err.Error())
	}
	
	if priceChangePercent >= 10 || priceChangePercent <= -5 {
		// todo
	}
}

/*
获取当前价格
获取当前货币的均价
*/

func (tradesUseCase *TradesUseCase) LastTicker(ctx context.Context) {
	price, err := market.LastPrice("BNBUSDT")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	fmt.Println("price: ", price)
}

/*
获取当前最佳挂单，参考该挂单进行成交
即: 买卖方之间价格博弈
Q:
	卖方价高，且量小，是否可以理解为涨？
	买方价低，且量小，是否理解为跌？ -- 没有买方需求？买房大量筹码位置？位置是否容易突破？
*/

func (tradesUseCase *TradesUseCase) GoodTicker(ctx context.Context) {
	goodTickerInfo, err := market.GoodTicker("BNBUSDT")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	fmt.Println("goodTickerInfo: ", goodTickerInfo)
}

func (tradesUseCase *TradesUseCase) ListSymbolPrice(ctx context.Context, symbol string, startTime int64, endTime int64) ([]TickerPrice, error) {
	return tradesUseCase.marketDbRepo.GetTickerPriceBySymbol(ctx, symbol, startTime, endTime)
}

/*
一段时间内价格信息 - 可以参考1分钟，2分钟，5分钟的价格变化

场景:
	每5分钟/每分钟计算价格变动，如果价格变动超过某个阀值，认为可以人工介入参考，类似
*/

func (tradesUseCase *TradesUseCase) TickerPrice(ctx context.Context, symbol string, windowsSize int64) {
	tickerPrice, err := market.TickerPrice(symbol, market.WindowsSize(windowsSize))
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息失败, err: %s", symbol, windowsSize, err.Error())
		return
	}
	
	priceChangePercent, err := strconv.ParseFloat(tickerPrice.PriceChangePercent, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换更改百分比失败(PriceChangePercent: %s), err: %s",
			symbol, windowsSize, tickerPrice.PriceChangePercent, err.Error())
		return
	}
	
	weightedAvgPrice, err := strconv.ParseFloat(tickerPrice.WeightedAvgPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(WeightedAvgPrice: %s), err: %s",
			symbol, windowsSize, tickerPrice.WeightedAvgPrice, err.Error())
		return
	}
	
	openPrice, err := strconv.ParseFloat(tickerPrice.OpenPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(OpenPrice: %s), err: %s",
			symbol, windowsSize, tickerPrice.OpenPrice, err.Error())
		return
	}
	
	highPrice, err := strconv.ParseFloat(tickerPrice.HighPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(HighPrice: %s), err: %s",
			symbol, windowsSize, tickerPrice.WeightedAvgPrice, err.Error())
		return
	}
	
	lowPrice, err := strconv.ParseFloat(tickerPrice.LowPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(LowPrice: %s), err: %s",
			symbol, windowsSize, tickerPrice.LowPrice, err.Error())
		return
	}
	
	lastPrice, err := strconv.ParseFloat(tickerPrice.LastPrice, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(LastPrice: %s), err: %s",
			symbol, windowsSize, tickerPrice.LastPrice, err.Error())
		return
	}
	
	volume, err := strconv.ParseFloat(tickerPrice.Volume, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(Volume: %s), err: %s",
			symbol, windowsSize, tickerPrice.Volume, err.Error())
		return
	}
	
	quoteVolume, err := strconv.ParseFloat(tickerPrice.QuoteVolume, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(QuoteVolume: %s), err: %s",
			symbol, windowsSize, tickerPrice.QuoteVolume, err.Error())
		return
	}
	
	priceChange, err := strconv.ParseFloat(tickerPrice.PriceChange, 32)
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 转换失败(PriceChange: %s), err: %s",
			symbol, windowsSize, tickerPrice.PriceChange, err.Error())
		return
	}
	
	if openPrice == 0 || highPrice == 0 {
		tradesUseCase.log.Errorf("获取%s (%d)秒内价格信息, 返回内容检查失败, openPrice: %f, openTime: %d",
			symbol, windowsSize, openPrice, tickerPrice.OpenTime/1000)
		return
	}
	
	err = tradesUseCase.marketDbRepo.RecordTickerPrice(ctx, TickerPrice{
		Symbol:             tickerPrice.Symbol,
		PriceChange:        priceChange,
		PriceChangePercent: priceChangePercent,
		WeightedAvgPrice:   weightedAvgPrice,
		OpenPrice:          openPrice,
		HighPrice:          highPrice,
		LowPrice:           lowPrice,
		LastPrice:          lastPrice,
		Volume:             volume,
		QuoteVolume:        quoteVolume,
		OpenTime:           tickerPrice.OpenTime / 1000,
		CloseTime:          tickerPrice.CloseTime / 1000,
		FirstId:            tickerPrice.FirstId,
		LastId:             tickerPrice.LastId,
		Count:              tickerPrice.Count,
	})
	
	if err != nil {
		tradesUseCase.log.Errorf("获取%s (%s)秒内价格信息, 记录数据到数据库失败, err: %s",
			symbol, windowsSize, err.Error())
	}
	
	// 如果价格波动超过 x，调用消息模块发送消息
	if priceChangePercent >= 10 || priceChangePercent <= -5 {
		// todo
	}
	
	return
}

/*
近几分钟平均单价 (5分钟)
*/

func (tradesUseCase *TradesUseCase) AvgPrice(ctx context.Context, symbol string) {
	avgPrice, err := market.AvgPrice(symbol)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	
	fmt.Println("avgPrice: ", avgPrice)
}
