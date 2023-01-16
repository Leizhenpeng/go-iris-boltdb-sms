package main

import (
	"github.com/kataras/iris/v12"
	"github.com/spf13/viper"
	"leizhenpeng/go-iris-boltdb-sms/controller"
	"leizhenpeng/go-iris-boltdb-sms/model"
	"leizhenpeng/go-iris-boltdb-sms/service"
)

func main() {
	//config.yaml load

	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	service.Init()
	if err != nil {
		panic(err)
		return
	}
	model.InitDb("sms.db")
	app := iris.Default()
	iris.RegisterOnInterrupt(func() {
		model.DbNow.Close()
	})
	registerController(app)
	app.Listen(":8080")

}

func registerController(app *iris.Application) {
	sms := app.Party("/sms")
	sms.Post("/new", controller.NewCode)
	sms.Post("/check", controller.CheckCode)
	sms.Get("/total", controller.Total)
	sms.Post("/clear", controller.Flush)
}
