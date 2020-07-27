package controllers

import (
	"VueAdmin/models"
	"VueAdmin/utils"
	_ "VueAdmin/utils"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"golang.org/x/crypto/bcrypt"
)

// Operations about Users
type UserController struct {
	BaseController
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	//uid := models.AddUser(user)
	//u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	var users []*models.User
	o := orm.NewOrm()
	o.QueryTable(new(models.User)).All(&users)
	u.Data["json"] = utils.APIRsp{Code: 20000, Msg: "ok", Data: users}
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /info [get]
func (u *UserController) Info() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.Data["json"] = utils.APIRsp{Code: 20000, Msg: "查询成功!", Data: map[string]string{"name": "running", "avator": "a"}}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		//uu, err := models.UpdateUser(uid, &user)
		//if err != nil {
		//	u.Data["json"] = err.Error()
		//} else {
		//	u.Data["json"] = uu
		//}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	//uid := u.GetString(":uid")
	//models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"用户名"
// @Param	password		query 	string	true		"密码"
// @Success 200 {json} {"code":200,"msg":"登陆成功!",data{"token":"aaa"}}
// @Failure 500 服务器错误
// @router /login [post]
func (u *UserController) Login() {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &data)
	if err != nil {
		u.Data["json"] = utils.APIRsp{Code: 40000, Msg: "密码或用户名不能为空!"}
		u.ServeJSON()
		return
	}
	// 校验参数
	valid := validation.Validation{}
	valid.Required(data.Username, "username")
	valid.Required(data.Password, "password")
	if valid.HasErrors() {
		u.Data["json"] = utils.APIRsp{Code: 40000, Msg: "密码或用户名不能为空!"}
		u.ServeJSON()
		return
	}
	user := models.User{}
	o := orm.NewOrm()
	err = o.QueryTable(user.TableName()).Filter("username", data.Username).One(&user)
	if err == orm.ErrNoRows {
		u.Data["json"] = utils.APIRsp{Code: 40000, Msg: "用户或密码错误!"}
		u.ServeJSON()
		return
	} else if err != nil {
		u.Data["json"] = utils.APIRsp{Code: 40000, Msg: "服务器错误!"}
		u.ServeJSON()
		return
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		u.Data["json"] = utils.APIRsp{Code: 40000, Msg: "用户或密码错误!"}
		u.ServeJSON()
		return
	}
	u.Data["json"] = utils.APIRsp{Code: 20000, Msg: "登陆成功!", Data: map[string]string{"token": models.CreateToken(&user)}}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [post]
func (u *UserController) Logout() {
	u.Data["json"] = utils.APIRsp{Code: 20000, Msg: "登出成功!"}
	u.ServeJSON()
}
