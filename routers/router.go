package routers

import (
	"Demo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/WS1", &controllers.CreateController{},"get:WS")
	beego.Router("/WS2", &controllers.ModifyController{},"get:WS")
	beego.Router("/WS3", &controllers.IssueController{},"get:WS")
	beego.Router("/WS4", &controllers.HistoricController{},"get:WS")
	beego.Router("/views/historic.html", &controllers.HistoricController{})
	beego.Router("/views/create.html", &controllers.CreateController{})
	beego.Router("/views/modify.html", &controllers.ModifyController{})
	beego.Router("/views/issue.html", &controllers.IssueController{})
	
}
