package service

import (
	"github.com/chainxx/bitx/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/startopsz/rule/pkg/response/errCode"
)

type ListSymbolPriceReq struct {
	Symbol    string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty" form:"symbol" validate:"required"`
	StartTime int64  `protobuf:"varint,2,opt,name=startTime,proto3" json:"startTime,omitempty" form:"startTime" validate:"required"`
	EndTime   int64  `protobuf:"varint,3,opt,name=endTime,proto3" json:"endTime,omitempty" form:"endTime" validate:"required"`
}

func (bitxService *BitxService) ListSymbolPrice(c *gin.Context) {
	req := &ListSymbolPriceReq{}
	err := common.BindUriQuery(c, req)
	if err != nil {
		return
	}
	
	tickerPrices, err := bitxService.trades.ListSymbolPrice(c.Request.Context(), req.Symbol, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(200, gin.H{"errCode": 500, "errMsg": "internal server error", "data": ""})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "data": tickerPrices})
	return
}
