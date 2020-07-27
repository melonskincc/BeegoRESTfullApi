package controllers

import (
	"VueAdmin/models"
	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
	User *models.User
}
