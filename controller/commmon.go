package controller

import "github.com/kataras/iris/v12"

type CommonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c CommonResponse) Send(ctx iris.Context) {
	ctx.JSON(c)
}

type ICommonResponse interface {
	Success() CommonResponse
	Fail() CommonResponse
	SetData(data interface{}) CommonResponse
	Send(ctx iris.Context)
}

var _ ICommonResponse = (*CommonResponse)(nil)

func (c CommonResponse) Success() CommonResponse {
	c.Code = 200
	c.Message = "success"
	return c
}

func (c CommonResponse) Fail() CommonResponse {
	c.Code = 500
	c.Message = "fail"
	return c
}

func (c CommonResponse) SetData(data interface{}) CommonResponse {
	c.Data = data
	return c
}
