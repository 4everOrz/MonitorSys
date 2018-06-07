package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["MonitorSys/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:AppController"],
		beego.ControllerComments{
			Method:           "AppLogin",
			Router:           `/alogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:AppController"],
		beego.ControllerComments{
			Method:           "AppRegist",
			Router:           `/aregist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:AppController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/getaall`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:AppController"],
		beego.ControllerComments{
			Method:           "UpdateApp",
			Router:           `/update`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})
	beego.GlobalControllerRouter["MonitorSys/controllers:AppController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:AppController"],
		beego.ControllerComments{
			Method:           "AppHeart",
			Router:           `/heart`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})
	beego.GlobalControllerRouter["MonitorSys/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:MdController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:MdController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           `/getall`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:MdController"],
		beego.ControllerComments{
			Method:           "GetAllFilter",
			Router:           `/getfilter`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:MdController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:MdController"],
		beego.ControllerComments{
			Method:           "GetData",
			Router:           `/getone`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:RoleController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:RoleController"],
		beego.ControllerComments{
			Method:           "Getall",
			Router:           `/rgetall`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"],
		beego.ControllerComments{
			Method:           "Addone",
			Router:           `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"],
		beego.ControllerComments{
			Method:           "GetServerAll",
			Router:           `/all`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"],
		beego.ControllerComments{
			Method:           "GetServerFilter",
			Router:           `/filter`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:ServerController"],
		beego.ControllerComments{
			Method:           "Updateone",
			Router:           `/update`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Changekey",
			Router:           `/changekey`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Getall",
			Router:           `/getall`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           `/ulogin`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Logout",
			Router:           `/ulogout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:UserController"],
		beego.ControllerComments{
			Method:           "Regist",
			Router:           `/uregist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:UserController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:UserController"],
		beego.ControllerComments{
			Method:           "UpdateOne",
			Router:           `/uupdate`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Params:           nil})

	beego.GlobalControllerRouter["MonitorSys/controllers:WebSocketController"] = append(beego.GlobalControllerRouter["MonitorSys/controllers:WebSocketController"],
		beego.ControllerComments{
			Method:           "Join",
			Router:           `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Params:           nil})

}
