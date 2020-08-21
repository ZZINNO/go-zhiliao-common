package network

import "github.com/gin-gonic/gin"

//通用方法

type ApiHandlerCommon struct{}

type ApiHandlerBaseResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetSuccessResp(data interface{}) ApiHandlerBaseResp {
	return ApiHandlerBaseResp{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func GetFailResp(err error, data interface{}) ApiHandlerBaseResp {
	return ApiHandlerBaseResp{
		Code:    1,
		Message: err.Error(),
		Data:    data,
	}
}

func (api ApiHandlerBaseResp) Write_resp(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"code": api.Code, "message": api.Message, "data": api.Data})
}

func (api *ApiHandlerCommon) Write_resp(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(200, gin.H{"code": code, "message": msg, "data": data})
}

func (api *ApiHandlerCommon) NoMethodResp(ctx *gin.Context) {
	ctx.JSON(405, gin.H{"code": 1, "message": "方法不被允许", "data": ""})
}
