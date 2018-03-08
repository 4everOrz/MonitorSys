package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"get","post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
