package main

import (
	"VueAdmin/models"
	_ "VueAdmin/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	//透明static
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
	beego.Run()
}

// 输入根目录转到login页面
func TransparentStatic(ctx *context.Context) {
	if strings.Index(ctx.Request.URL.Path, "v1/") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "dist/"+ctx.Request.URL.Path)
}

func init() {
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	_ = orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("SqlConn"))
	debug, _ := beego.AppConfig.Bool("debug")
	orm.Debug = debug
	orm.SetMaxIdleConns("default", 10)
	orm.SetMaxOpenConns("default", 50)
	_ = orm.RunSyncdb("default", false, true)
	encodePwd, _ := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	u := models.User{Username: "running", Password: string(encodePwd), NickName: "cia"}
	o := orm.NewOrm()
	_, _ = o.Insert(&u)
}
