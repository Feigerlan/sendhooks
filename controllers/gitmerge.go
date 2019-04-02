package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"reflect"
	"sendhooks/models"
)


// Operations about object
type GitmergeController struct {
	beego.Controller
}
//定义接收到的合并请求结构体
type MergeRequest struct {
	Object_kind string 						`json:"object_kind"`
	User User 								`json:"user"`
	Project Project							`json:"project"`
	Object_attributes Object_Attributes     `json:"object_attributes"`
	Labeles bool                            `json:"labeles"`
	Repository Repository                   `json:"repository"`
}
//定义接收到的合并请求结构体中user
type User struct {
	Name string           `json:"name"`
	Username string       `json:"username"`
	Avatar_url string     `json:"avatar_url"`
}
//定义接收到的合并请求结构体中的project
type Project struct {
	 Name string           `json:"name"`
	 Description string    `json:"description"`
	 Web_url string 		`json:"Web_url"`
	 Git_ssh_url string     `json:"git_ssh_url"`
	 Git_http_url string `json:"git_http_url"`
}

type Object_Attributes struct {
	Id int 								`json:"id"`
	Target_branch string				`json:"target_branch"`
	Title string  						`json:"title"`
	Created_at string  					`json:"created_at"`
	Updated_at string 					`json:"updated_at"`
	Merge_status string					`json:"merge_status"`
}
//定义接收到的合并请求结构体Repository信息
type Repository struct {
	Name string 					`json:"name"`
	Url string                      `json:"url"`
	Description string 				`json:"description"`
	Homepage string 				`json:"homepage"`

}



func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (o *GitmergeController) Post() {
	var ob MergeRequest
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	//fmt.Printf(ob.Object_kind)
	if ob.Object_kind == "merge_request" {
           fmt.Println("接收到合并请求！")
           //fmt.Printf("%+v",ob)
           fmt.Printf(typeof(ob))
		   o.Data["json"] = ob
		   o.ServeJSON()
		   sendmsg("http://baidu.com")
		} else{
		fmt.Println("请求参数错误！")
		o.Ctx.Output.Status = 402
		o.Data["json"] = "请求参数错误！"
		o.ServeJSON()
	}

}


func sendmsg(url string){
	var mm = make(map[string]interface{})
	mm["userID"] = "lanxiahui"
	mm["pwd"] = "123456"
	jsonStr, err := json.Marshal(mm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//初始化一个http客户端
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	//取出body
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *GitmergeController) Get() {
	objectId := o.Ctx.Input.Param(":objectId")
	println(objectId)
	if objectId != "" {
		ob, err := models.GetOne(objectId)
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = ob
		}
	}
	o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *GitmergeController) GetAll() {
	println("123")
	obs := models.GetAll()
	o.Data["json"] = obs
	o.ServeJSON()
}

// @Title Update
// @Description update the object
// @Param	objectId		path 	string	true		"The objectid you want to update"
// @Param	body		body 	models.Object	true		"The body"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [put]
func (o *GitmergeController) Put() {
	objectId := o.Ctx.Input.Param(":objectId")
	var ob models.Object
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)

	err := models.Update(objectId, ob.Score)
	if err != nil {
		o.Data["json"] = err.Error()
	} else {
		o.Data["json"] = "update success!"
	}
	o.ServeJSON()
}

// @Title Delete
// @Description delete the object
// @Param	objectId		path 	string	true		"The objectId you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 objectId is empty
// @router /:objectId [delete]
func (o *GitmergeController) Delete() {
	objectId := o.Ctx.Input.Param(":objectId")
	models.Delete(objectId)
	o.Data["json"] = "delete success!"
	o.ServeJSON()
}

