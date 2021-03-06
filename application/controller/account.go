package controller

import (
	"bytes"
	"net/http"

	bc "gdzce.cn/perishable-food/application/blockchain"
	"github.com/gin-gonic/gin"
)

// 账户查询请求体
type accountListRequestBody struct {
	AccountId string `form:"account_id" json:"account_id" binding:"required"`
}

// 查询账户列表
func AccountList(ctx *gin.Context) {
	// 解析请求体
	req := new(accountListRequestBody)
	if err := ctx.ShouldBind(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 调用链码的queryAccount，查询账户列表
	resp, err := bc.ChannelQuery("queryAccount", [][]byte{
		[]byte(req.AccountId),
	})
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 将结果返回
	ctx.String(http.StatusOK, bytes.NewBuffer(resp.Payload).String())
}
