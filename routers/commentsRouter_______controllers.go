package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:AppController"],
		beego.ControllerComments{
			Method: "AppLogin",
			Router: `/alogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:AppController"],
		beego.ControllerComments{
			Method: "AppRegist",
			Router: `/aregist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/getall`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:MdController"],
		beego.ControllerComments{
			Method: "GetData",
			Router: `/getone`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:ServerController"],
		beego.ControllerComments{
			Method: "Addone",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:ServerController"],
		beego.ControllerComments{
			Method: "GetServerAll",
			Router: `/all`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:ServerController"],
		beego.ControllerComments{
			Method: "GetServerFilter",
			Router: `/filter`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Changekey",
			Router: `/changekey`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Getall",
			Router: `/getall`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/ulogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/ulogout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:UserController"],
		beego.ControllerComments{
			Method: "Regist",
			Router: `/uregist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:WebSocketController"] = append(beego.GlobalControllerRouter["MonitoringSystemAPI/controllers:WebSocketController"],
		beego.ControllerComments{
			Method: "Join",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
