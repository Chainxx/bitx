package service

import (
	"github.com/chainxx/bitx/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/startopsz/rule/pkg/response/errCode"
)

type GetDoubleColorBallCountReq struct {
	RedBallCount  int `protobuf:"varint,2,opt,name=startTime,proto3" json:"redBallCount,omitempty" form:"redBallCount" validate:"required"`
	BlueBallCount int `protobuf:"varint,2,opt,name=startTime,proto3" json:"blueBallCount,omitempty" form:"blueBallCount" validate:"required"`
}

func (bitxService *BitxService) GetDoubleColorBallCount(c *gin.Context) {
	req := &GetDoubleColorBallCountReq{}
	err := common.BindUriQuery(c, req)
	if err != nil {
		return
	}
	
	betting := bitxService.lotto.GetDoubleColorBallCount(c.Request.Context(), req.RedBallCount, req.BlueBallCount)
	
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "data": betting})
	return
}
