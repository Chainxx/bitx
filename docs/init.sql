create database bitx charset = utf8mb4;
create user 'bit'@'%' identified by '083b594X#9e51ad3';
grant all on bitx.* to 'bit'@'%';
flush privileges;


drop table if exists daily_asset;
create table if not exists daily_asset
(
    asset       varchar(100) not null comment '资产名',
    free        double default 0 comment '可操作资产',
    locked      double default 0 comment '锁定资产',
    time        bigint comment '时间',
    create_time bigint comment '创建时间',
    key idx_time (time),
    unique asset_time (asset, time)
) comment '资产表';


drop table if exists recharge;
create table if not exists recharge
(
    id             varchar(50)  not null primary key comment 'id',
    amount         double default 0 comment '充值值',
    coin           varchar(100) not null comment '币名',
    network        varchar(100) not null comment '网络',
    status         int comment '状态',
    address        varchar(200) not null comment '地址',
    address_tag    varchar(100) not null comment '地址标签',
    tx_id          varchar(200) not null comment '交易id',
    insert_time    bigint,
    transfer_type  int,
    confirm_times  varchar(100),
    unlock_confirm int,
    wallet_type    int
) comment '充值表';



drop table if exists ticker_price;
create table if not exists ticker_price
(
    symbol               varchar(100) not null comment '代号',
    price_change         double default 0 comment '价格变化',
    price_change_percent double default 0 comment '价格变化百分比',
    weighted_avg_price   double default 0 comment '加权平均价格',
    open_price           double default 0 comment '开盘价, 开始价格',
    high_price           double default 0 comment '高价',
    low_price            double default 0 comment '低价',
    last_price           double default 0 comment '最后价格',
    volume               double default 0 comment '体量',
    quote_volume         double default 0 comment '报价货币成交量',
    open_time            bigint comment 'ticker的开始时间',
    close_time           bigint comment 'ticker的结束时间',
    first_id             bigint comment '统计时间内的第一笔trade id',
    last_id              bigint comment '统计时间内的最后一笔trade id',
    count                bigint comment '统计时间内交易笔数'
) comment '价格表';

drop table if exists ticker_price_daily;
create table if not exists ticker_price_daily
(
    symbol               varchar(100) not null comment '代号',
    price_change         double default 0 comment '价格变化',
    price_change_percent double default 0 comment '价格变化百分比',
    weighted_avg_price   double default 0 comment '加权平均价格',
    prev_close_price     double default 0 comment '前一交易周期的收盘价格',
    bid_price            double default 0 comment '出价价格',
    bid_qty              double default 0 comment '出价数量',
    ask_price            double default 0 comment '询价价格',
    ask_qty              double default 0 comment '询价数量',
    open_price           double default 0 comment '开盘价, 开始价格',
    high_price           double default 0 comment '高价',
    low_price            double default 0 comment '低价',
    last_price           double default 0 comment '最后价格',
    last_qty             double default 0 comment '最后价格数量',
    volume               double default 0 comment '体量',
    quote_volume         double default 0 comment '报价货币成交量',
    open_time            bigint comment 'ticker的开始时间',
    close_time           bigint comment 'ticker的结束时间',
    first_id             bigint comment '统计时间内的第一笔trade id',
    last_id              bigint comment '统计时间内的最后一笔trade id',
    count                bigint comment '统计时间内交易笔数'
) comment '价格表';



drop table if exists transport;
create table if not exists transport
(
    id                varchar(150) not null comment '该笔提现在币安的id',
    amount            double default 0 comment '提现转出金额',
    transaction_fee   double default 0 comment '手续费',
    coin              varchar(150) not null comment '',
    status            int    default 0 comment '',
    address           varchar(150) not null comment '',
    tx_id             varchar(150) not null comment '提现交易id',
    apply_time        varchar(150) comment '',
    network           varchar(150) comment '',
    transfer_type     varchar(150) comment '1: 站内转账, 0: 站外转账',
    withdraw_order_id varchar(150) comment '自定义ID, 如果没有则不返回该字段',
    info              text comment '提币失败原因',
    confirm_no        int comment '提现确认数',
    wallet_type       int comment '1: 资金钱包 0:现货钱包',
    tx_key            varchar(150) comment '',
    complete_time     varchar(150) comment '提现完成，成功下账时间(UTC)'
) comment '提币记录';


