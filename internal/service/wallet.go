package service

import (
	"github.com/chainxx/bitx/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/startopsz/rule/pkg/response/errCode"
)

type GetDailyAssetReq struct {
	Asset     string `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset,omitempty" form:"asset" validate:"required"`
	StartTime int64  `protobuf:"varint,2,opt,name=startTime,proto3" json:"startTime,omitempty" form:"startTime" validate:"required"`
	EndTime   int64  `protobuf:"varint,3,opt,name=endTime,proto3" json:"endTime,omitempty" form:"endTime" validate:"required"`
}

func (bitxService *BitxService) GetDailyAsset(c *gin.Context) {
	req := &GetDailyAssetReq{}
	err := common.BindUriQuery(c, req)
	if err != nil {
		return
	}
	
	dailyAssets, err := bitxService.wallet.GetPerDayAssetsByDB(c.Request.Context(), req.Asset, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(200, gin.H{"errCode": 500, "errMsg": "internal server error", "data": ""})
		return
	}
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "data": dailyAssets})
	return
}
