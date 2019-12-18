package controllers

import "github.com/astaxie/beego"
import (
		"time"
		"strconv"
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
	projectid = make(chan int, 10)
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
		case version := <- issue:
			tag := <- tags
			projectID := <- projectid
			client := <- ch3
			tagname := tag.Name
		
			beego.Info("更新标签名" + tagname)
			tagname = VersionNumberUpdate(tagname,version.IssueStatus)
			beego.Info("更新标签名后" + tagname)
			
			mesStr, _ := models.HTTPPostTag(projectID, tagname)
			//利用cnt判断是否是最后一行数据发布
			cnt := 0
			for i := 0; i < len(version.ServiceList); i++ {
				if version.ServiceList[i].ServiceNumber == "&&&" {
					cnt++
				}
				if version.ServiceList[i].ServiceNumber == strconv.Itoa(projectID) {
					version.ServiceList[i].ServiceNumber = "&&&"
				}
			}
			beego.Info(cnt)
			beego.Info(len(version.ServiceList))
			if cnt == len(version.ServiceList) - 1 {
				version.IssueTime = time.Now()
				version.Status = true
				services := GetGroups()
				var mp = make(map[string]string)
				for _, val := range services {
					mp[val.ServiceName] = val.ServiceNumber
				}
				
				for i := 0; i < len(version.ServiceList); i++ {
					version.ServiceList[i].ServiceNumber = mp[version.ServiceList[i].ServiceName]
				}
				models.IssueVersion(&version)
			} else {
				models.UpdateVersion(&version)
			}
			//将创建标签的返回消息，返回到前端，用消息框弹出
			data, err := json.Marshal(mesStr)
			checkErr(err)
			client.conn.WriteMessage(websocket.TextMessage, data)
		case merges := <- merge:
			client := <- ch3
			//多个合并请求时，一一给与合并，其中一个出问题，便停止并将失败原因传送到前端
			for _, MR := range merges {
				beego.Info(MR.ProjectID)
				beego.Info(MR.IID)
				mergemessage := models.Merge(MR.ProjectID, MR.IID)
				if mergemessage != "" {
					beego.Info(mergemessage)
					data, err := json.Marshal(mergemessage)
					checkErr(err)
					client.conn.WriteMessage(websocket.TextMessage, data)
					break
				}
			}
		}

	}
}