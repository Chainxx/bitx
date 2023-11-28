# bitx

bitx 是一个定时获取 币安API 数据的定时程序 & HTTP REST API 程序

## 背景

在玩数据货币时，想通过 币安API 获取数据到本地，然后根据自己的需求制图查看走势。

bitx 配置好对应的 APIKey 之后，就可以定时获取数据到自己的 MySQL 数据库中。

## 量化交易/高频交易

根据 币安API & 币安API数据 & 本地历史数据，编写自己的算法，进行量化交易

### Example

1. 计算交易金额，手续费金额，波动比例盈利
2. 波动交易笔数，补仓行为？

计算阶梯上行/下行的可行性，计算买点和卖点的获利行为

动态计算获利点？比如通过震荡计算趋势获利百分比

波动点数据:

1. 昨日最低点
2. 昨日收盘点  (24小时运行制下有意义么？)

## Usage

```shell
source env.properties  && ./bin/task
```

## 制图

原计划是编写一个简单的web页面，通过rest api获取数据。

推荐方案:

    使用 Grafana 添加 MySQL 数据源，根据自己的需求，编写 SQL 制图。

    可以参考倒入 docs/grafana.json 配置

## 部署

部署服务需要部署在可以访问 币安API 的地区，请根据需求修改 env.properties 文件

## 关联文献

[币安API](https://binance-docs.github.io/apidocs/spot/cn/#4175e32579)

[ImToken](https://docs.token.im/tokenlon-open-api/?locale=zh)
