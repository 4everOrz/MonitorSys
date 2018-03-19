// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"MonitoringSystemAPI/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/app",
			beego.NSInclude(
				&controllers.AppController{},
			),
		), /**/
		beego.NSNamespace("/data",
			beego.NSInclude(
				&controllers.MdController{},
			),
		),
		beego.NSNamespace("/ws",
			beego.NSInclude(
				&controllers.WebSocketController{},
			),
		),
		beego.NSNamespace("/sv",
			beego.NSInclude(
				&controllers.ServerController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
