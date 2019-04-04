// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"sendhooks/controllers"
)

//token
var (
	key []byte = []byte("xwfintech")
)

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
//定义响应json 结构体
type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
}


func init() {
	//生成token
	token := GenToken()
	fmt.Println(token)
	//token鉴权过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {
		//fmt.Printf("%+v",ctx.Input.Context.Request)
		authString := ctx.Input.Header("X-Gitlab-Token")
		fmt.Println("接收到带来的token:" + authString)
		if !CheckToken(authString){   //如果检查token不通过
			ctx.Output.Status = 406    //设置返回码
			//设置返回字符串
			resstr := `{"data":{"code":406 ,"msg":"token不合法"}}`
			//初始化返回信息
			var res Response
			if err := json.Unmarshal([]byte(resstr), &res); err == nil {
				//反序列化为json，如果无错则打印
				fmt.Println(res)
			} else {
				fmt.Println(err)
			}
	        //返回给客户端信息
			ctx.Output.JSON(res,false,false)
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
	fmt.Println(token)
	_, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println("parase with claims failed.", err)
		return false
	}
	return true
}

