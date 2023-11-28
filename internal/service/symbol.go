package service

import (
	"github.com/gin-gonic/gin"
	"github.com/startopsz/rule/pkg/response/errCode"
)

func (bitxService *BitxService) ListSymbol(c *gin.Context) {
	symbols := bitxService.trades.ListSymbols(c.Request.Context())
	c.JSON(200, gin.H{"errCode": errCode.NormalCode, "errMsg": errCode.NormalMsg, "data": symbols})
	return
}
