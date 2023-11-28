# 币安API

isBuyerMaker = true -> SELL
isBuyerMaker = false -> BUY

## api地址

https://api.binance.com

https://api1.binance.com

https://api2.binance.com

https://api3.binance.com

## 是intervalLetter

intervalLetter 作为头部值:

    SECOND => S
    MINUTE => M
    HOUR => H
    DAY => D

## Code

收到429后仍然继续违反访问限制，会被封禁IP，并收到418错误码

## API-keys

如果需要 API-keys，应当在HTTP头中以 X-MBX-APIKEY字段传递。
API-keys 与 secret-keys 是大小写敏感的。

NONE	        不需要鉴权的接口
TRADE	        需要有效的 API-Key 和签名
MARGIN	        需要有效的 API-Key 和签名
USER_DATA	    需要有效的 API-Key 和签名
USER_STREAM	    需要有效的 API-Key
MARKET_DATA	    需要有效的 API-Key

-H "X-MBX-APIKEY: api-key"

## 交易对状态 (状态 status)

    PRE_TRADING 交易前
    TRADING 交易中
    POST_TRADING 交易后
    END_OF_DAY
    HALT
    AUCTION_MATCH
    BREAK

## 交易对类型

    SPOT 现货
    MARGIN 杠杆
    LEVERAGED 杠杆代币
    TRD_GRP_002 交易组 002
    TRD_GRP_003 交易组 003
    TRD_GRP_004 交易组 004

## 参考

[API 开发文档](https://binance-docs.github.io/apidocs/spot/cn/#b122f813d5)

[币安 API 导航页](https://www.binance.com/zh-CN/binance-api)

[如何创建API](https://www.binance.com/zh-CN/support/faq/360002502072)
