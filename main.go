package main

import (
	"github.com/astaxie/beego"
	"moshopserver/models"
	_ "moshopserver/routers"
	_ "moshopserver/utils"
)

func init() {
	models.InitDB()
}

func main() {

	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.CopyRequestBody = true

	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.Listen.HTTPAddr = "192.168.10.128"
	beego.BConfig.Listen.HTTPPort = 8080

	//beego.InsertFilter("/api/*", beego.BeforeExec, services.FilterFunc, true, true)

	beego.Run() // listen and serve on 0.0.0.0:8080

}
