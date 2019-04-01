package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"sendhooks/models"
)

type GiteventController struct {
	beego.Controller
}

type Events struct {
	Event_name 	string 				`json:"event_name"`
	Project Project							`json:"project"`
	project_id int                          `json:"project_id"`
	Repository Repository                   `json:"repository"`
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (o *GiteventController) Post() {
	var ob Events
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	//fmt.Printf("%+v",ob)
	//fmt.Printf(ob.Object_kind)
	if ob.Event_name == "tag_push" {
           fmt.Println("接收到tag请求！")
           o.Data["json"] = ob
		   o.ServeJSON()
		}
	if ob.Event_name == "repository_update" {
		fmt.Println("接收到仓库更新请求！")
		o.Data["json"] = ob
		o.ServeJSON()
	}
	if ob.Event_name == "push" {
		fmt.Println("接收到puth请求！")
		o.Data["json"] = ob
		o.ServeJSON()
	}else{
		fmt.Println("请求参数错误！")
		o.Ctx.Output.Status = 402
		o.Data["json"] = "请求参数错误！"
		o.ServeJSON()
	}

}

//func (o *GitlabController) Post() {
//	var merge_request string
//	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
//	ObjectId := models.AddOne(ob)
//	o.Data["json"] = map[string]string{"ObjectId": objectid}
//	o.ServeJSON()
//}

// @Title Get
// @Description find object by objectid
// @Param	objectId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:objectId [get]
func (o *GiteventController) Get() {
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
func (o *GiteventController) GetAll() {
	//println("123")
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
func (o *GiteventController) Put() {
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
func (o *GiteventController) Delete() {
	objectId := o.Ctx.Input.Param(":objectId")
	models.Delete(objectId)
	o.Data["json"] = "delete success!"
	o.ServeJSON()
}

