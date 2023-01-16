package controller

import (
	"github.com/spf13/viper"
	"leizhenpeng/go-iris-boltdb-sms/service"

	"github.com/kataras/iris/v12"
)

type NewCodeRequest struct {
	Phone string `json:"phone"`
}

func NewCode(ctx iris.Context) {
	var req NewCodeRequest
	err := ctx.ReadJSON(&req)
	if err != nil {
		return
	}
	phone := req.Phone
	if !service.Sms.ValidPhone(phone) {
		CommonResponse{}.Fail().SetData("not valid phone").Send(ctx)
		return
	}
	if service.Sms.IfExist(phone) {
		CommonResponse{}.Fail().SetData(
			iris.Map{
				"phone": phone,
				"msg":   "already send code",
			}).Send(ctx)
		return
	}
	code, err := service.Sms.GenCode(phone)
	if err != nil {
		return
	}
	CommonResponse{}.Success().SetData(iris.Map{
		"phone":  phone,
		"code":   code,
		"expire": viper.GetString("EXPIRE") + "s",
	}).Send(ctx)

}

type CheckCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func CheckCode(ctx iris.Context) {
	var req CheckCodeRequest
	err := ctx.ReadJSON(&req)
	if err != nil {
		return
	}
	if service.Sms.ValidCode(req.Phone, req.Code) {
		CommonResponse{}.Success().SetData("pass").Send(ctx)
		return
	}
	CommonResponse{}.Fail().SetData("check phone code fail").
		Send(ctx)
}

func Total(ctx iris.Context) {
	CommonResponse{}.Success().SetData(iris.Map{
		"total": service.Sms.Total(),
	}).Send(ctx)
}

func Flush(ctx iris.Context) {
	service.Sms.ClearAll()
	CommonResponse{}.Success().SetData("flush success").Send(ctx)
}
