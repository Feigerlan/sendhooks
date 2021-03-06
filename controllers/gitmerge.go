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
//定义合并请求中的对象属性
type Object_Attributes struct {
	Id int 								`json:"id"`
	Target_branch string				`json:"target_branch"`
	Source_branch string				`json:"source_branch"`
	Title string  						`json:"title"`
	Created_at string  					`json:"created_at"`
	Updated_at string 					`json:"updated_at"`
	Merge_status string					`json:"merge_status"`
	Last_commit Last_commit				`json:"last_commit"`

}
//
type Last_commit struct {
	Id string		`json:"id"`
	Message string  `json:"message"`
	Timestamp string `json:"timestamp"`
	Url string `json:"url"`
    Author Author `json:"author"`
}
type Author struct {
	Name string `json:"name"`
	Email string `json:"email"`
}
//定义接收到的合并请求结构体Repository信息
type Repository struct {
	Name string 					`json:"name"`
	Url string                      `json:"url"`
	Description string 				`json:"description"`
	Homepage string 				`json:"homepage"`

}

//定义消息接口json结构体
type Messages struct {
	Type string     `json:"type"`
	Title string  	 `json:"title"`
	Content string   `json:"content"`
	Ways string      `json:"ways"`
	Receiver string  `json:"receiver"`
}




//获取数据类型
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
	//初始化一个合并请求对像用于接收hooks发过来的json
	var ob MergeRequest
	//将hooks发过来的json反序列化并交内容写到ob对像的地址中
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	//fmt.Printf(ob.Object_kind)
	//如果检测到其中有合并请求的字段
	if ob.Object_kind == "merge_request" {
           fmt.Println("接收到合并请求！")
           fmt.Printf("%+v",ob)

           //fmt.Printf(typeof(ob))
		   //o.Data["json"] = ob
		   //o.ServeJSON()
		   //调度发送消息
		   sendmsg("http://122.152.209.199:2046/api/v1/atlassian/message/",ob)

		   getowner()
		} else{        //如果没有合并请求的字段返回错误码
		fmt.Println("请求参数错误！")
		o.Ctx.Output.Status = 402
		o.Data["json"] = "请求参数错误！"
		o.ServeJSON()
	}

}


func sendmsg(url string,mm MergeRequest){
	proname := "来自"+ mm.Repository.Name + "项目的合并请求"
	var mess Messages
	mess.Type = "gitlab"
	mess.Title = proname + ":" + mm.Repository.Name
	mess.Ways = "wx"
	mess.Receiver = "feigerlan@xwfintech.com"
	mess.Content = "项目地址：" + mm.Project.Web_url + "\n-----------发起者：" + mm.Object_attributes.Last_commit.Author.Name
	jsonStr, err := json.Marshal(mess)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(jsonStr)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	//初始化一个http客户端
	client := &http.Client{}
	//发送请求，接收结果和错误
	resp, err := client.Do(req)
	//如果错误不为空，打印错误
	if err != nil {
		panic(err)
	}

	//最终关闭
	defer resp.Body.Close()
	//打印服务返回的状态码
	fmt.Println("响应状态码:", resp.Status)
	//打印返回的头部
	fmt.Println("响应头部:", resp.Header)
	//取出body
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("响应消息体:", string(body))
	if resp.Status == "200 OK"{
		fmt.Println("成功")
	}
}

func getowner(){
	pro, err := http.Get("http://gitlab.xwfintech.com/api/v4/projects/292?private_token=zC4YkdmxUr1_jBH9xy1x")
	if err != nil{
		println(err)
	}
	body, err := ioutil.ReadAll(pro.Body)
	defer pro.Body.Close()
	var proj map[string]interface{}
	proj = make(map[string]interface{})
	json.Unmarshal(body,&proj)
	fmt.Printf("%+v",proj)

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

