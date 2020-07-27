// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"VueAdmin/controllers"
	"VueAdmin/models"
	"VueAdmin/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
)

// jwt 鉴权
var TokenFilter = func(ctx *context.Context) {
	if ctx.Request.RequestURI != "/api/v1/user/login" {
		if ctx.Request.Header.Get("Authorization") == "" {
			ctx.Output.JSON(utils.APIRsp{Code: 50008, Msg: "未登陆！"}, false, false)
		}
		username, err := models.VerifyToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.Output.JSON(utils.APIRsp{Code: 50014, Msg: "身份信息有误,请重新登陆！"}, false, false)
		}
		ctx.Input.SetParam("username", username)
	}
}

func init() {
	// 定义路由
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	// 跨域请求设置                  #执行完了controller再执行
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	// token 校验				# 在执行controller前执行
	beego.InsertFilter("*", beego.BeforeExec, TokenFilter)
	//此处上面的是要添加的字段, 添加路由
	beego.AddNamespace(ns)
}
