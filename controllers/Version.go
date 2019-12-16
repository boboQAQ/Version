package controllers

import "github.com/astaxie/beego"
import (
		"time"
		_ "strconv"
       "Demo/models"
	   "encoding/json"
	   "github.com/gorilla/websocket"
)
// VersionController 版本控制器
type VersionController struct {
	beego.Controller
}

func checkErr (err error) {
	if err != nil {
		panic(err)
	}
}
// GetAllStatus 获取所有状态为未发布的版本
// @Description 返回所有的未发布版本的数据
// @Success 200 {object} models.Version
// @router / [get]
func (v *VersionController) GetAllStatus() {
	vs := models.GetAllStatus()
	vss, err := json.Marshal(vs)
	checkErr(err)

	v.Data["json"] = vss
	v.ServeJSON()
}
// GetByID 获得一个版本信息
// @Description 返回某未发布版本的数据
// @Param      id            path   int    true          "The key for staticblock"
// @Success 200 {object} models.Version
// @router /:id [get]
func (v *VersionController) GetByID() {
	id, _ := v.GetInt("id")

	var ver models.Version
	ver = models.GetByID(id)

	vers, err := json.Marshal(ver)
	checkErr(err)
	v.Data["json"] = vers

	v.ServeJSON()

}
// PostCreatVersion 创建一个新版本
func (v *VersionController) PostCreatVersion() {
	var version models.Version
	json.Unmarshal(v.Ctx.Input.RequestBody, &version)

	version.CreatTime = time.Now()
	version.Status = false
	id := models.CreatVersion(&version)
	v.Data["json"] = id;
	v.ServeJSON()
}

// Update issue小版本
func (v *VersionController) Update() {
	var version models.Version
	json.Unmarshal(v.Ctx.Input.RequestBody, &version)

	version.IssueTime = time.Now()
	models.UpdateVersion(&version)
}

var (
	create = make(chan models.Version, 10)
	update = make(chan models.Version, 10)
	issue  = make(chan models.Version, 10)
	tags    = make(chan models.Tags, 10)
	merge  = make(chan []models.MergeRequest, 10)

	//发布页面后端对前端的通信通道
	ch3 = make(chan Version, 1)

)
func init() {
	go  broadcaster()
}

func  broadcaster() {
	for {
		select {
		case version := <- create:
			
			version.CreatTime = time.Now()
			version.IssueTime = time.Now()
			version.Status = false
			_ = models.CreatVersion(&version)
		case version := <- update:
			version.IssueTime = time.Now()
			//beego.Info( VersionNumberUpdate(version.VersionNumber, 0) )
			//beego.Info( ServiceListUpdate(version.ServiceList, 0) )
			models.UpdateVersion(&version)
		//case tag := <- issue:
			//tagname := tag.Name
		
			//beego.Info("更新标签名" + tagname)
			// version.Status = true
			// version.IssueTime = time.Now()
			// beego.Info( VersionNumberUpdate(version.VersionNumber, version.IssueStatus) )
			// beego.Info( ServiceListUpdate(version.ServiceList, version.IssueStatus) )
			// models.IssueVersion(&version)
		case merges := <- merge:
			client := <- ch3
			for _, MR := range merges {
				mergemessage := models.Merge(MR.ID, MR.IID)
				if mergemessage != "" {
					data, err := json.Marshal(mergemessage)
					checkErr(err)
					client.conn.WriteMessage(websocket.TextMessage, data)
					break
				}
			}
		}

	}
}