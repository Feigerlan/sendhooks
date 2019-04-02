// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"sendhooks/controllers"
)

//token
var (
	key []byte = []byte("xwfintech")
)

func init() {
	//生成token
	token := GenToken()
	fmt.Println(token)
	//token鉴权过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {
		authString := ctx.Input.Header("Authorization")
		//fmt.Println("接收到带来的token:" + authString)
		if !CheckToken(authString){
			ctx.Output.Status = 404
			ctx.Output.JSON("{无效token}",false,false)
		}
		})

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/gitmerge",
			beego.NSInclude(
				&controllers.GitmergeController{},
			),
		),
		beego.NSNamespace("/gitevent",
			beego.NSInclude(
				&controllers.GiteventController{},
			),
		),
	)
	beego.AddNamespace(ns)
}






// 产生json web token
func GenToken() string {
	claims := &jwt.StandardClaims{
		//NotBefore: int64(time.Now().Unix()),
		//ExpiresAt: int64(time.Now().Unix() + 1000),
		Issuer:    "xwfintech",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return ss
}


// 校验token是否有效
func CheckToken(token string) bool {
	//将TOKEN解析
	//fmt.Println(token)
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return false
	}
	return true
}

