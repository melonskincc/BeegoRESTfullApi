package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id         int
	Username   string `orm:"unique"`
	Password   string
	NickName   string
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func (u *User) TableName() string {
	return "app_user"
}

//func AddUser(u User) string {
//u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
//UserList[u.Id] = &u
//return u.Id
//}
// 创建token
func CreateToken(u *User) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	exp, _ := strconv.Atoi(beego.AppConfig.String("TokenExp"))
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(exp)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = u.Username
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(beego.AppConfig.String("TokenSecret")))
	return tokenString
}

// 校验token
func VerifyToken(encodeToken string) (string, error) {
	username := ""
	s := strings.Split(encodeToken, " ")
	if len(s) != 2 || s[0] != "Bearer" {
		return username, errors.New("认证失败")
	}
	encodeToken = s[1]
	token, _ := jwt.Parse(encodeToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return username, errors.New("认证失败")
		}
		return []byte(beego.AppConfig.String("TokenSecret")), nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	err := claims.Valid()
	if err != nil {
		return username, errors.New("认证失败")
	}
	username = claims["username"].(string)
	return username, nil
}

func GetUser(uid string) (u *User, err error) {
	//if u, ok := UserList[uid]; ok {
	//	return u, nil
	//}
	return nil, errors.New("User not exists")
}

//func GetAllUsers() map[string]*User {
//return UserLis
//}

//func UpdateUser(uid string, uu *User) (a *User, err error) {
//	if u, ok := UserList[uid]; ok {
//		if uu.Username != "" {
//			u.Username = uu.Username
//		}
//		if uu.Password != "" {
//			u.Password = uu.Password
//		}
//		if uu.Profile.Age != 0 {
//			u.Profile.Age = uu.Profile.Age
//		}
//		if uu.Profile.Address != "" {
//			u.Profile.Address = uu.Profile.Address
//		}
//		if uu.Profile.Gender != "" {
//			u.Profile.Gender = uu.Profile.Gender
//		}
//		if uu.Profile.Email != "" {
//			u.Profile.Email = uu.Profile.Email
//		}
//		return u, nil
//	}
//	return nil, errors.New("User Not Exist")
//}
//
//func Login(username, password string) bool {
//	for _, u := range UserList {
//		if u.Username == username && u.Password == password {
//			return true
//		}
//	}
//	return false
//}
//
//func DeleteUser(uid string) {
//	delete(UserList, uid)
//}

func init() {
	orm.RegisterModel(new(User))
}
