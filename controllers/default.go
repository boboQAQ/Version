package controllers

import (
	"net/http"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"strconv"
	"Demo/models"
	"encoding/json"
)

type message struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

type sendmessage struct {
	Services []models.Service `json:"services"`
	Versions []models.Version `json:"versions"`
}
// MainController 默认控制器
type MainController struct {
	beego.Controller
}
//CreateController 创造控制器
type CreateController struct {
	beego.Controller
}
// ModifyController 修改控制器
type ModifyController struct {
	beego.Controller
}
// IssueController 发布控制器
type IssueController struct {
	beego.Controller
}

// HistoricController 发布控制器
type HistoricController struct {
	beego.Controller
}

// Get 菜单栏默认调用
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"

	c.TplName = "index.html"
	
	
	
}
// Get 创造页面默认调用
func (c *CreateController) Get() {
	c.TplName = "create.html"
}
// Get 修改页面默认调用
func (c *ModifyController) Get() {
	c.TplName = "modify.html"
}
// Get 发布页面默认调用
func (c *IssueController) Get() {
	c.TplName = "issue.html"
}
// Get 历史页面默认调用
func (c *HistoricController) Get() {
	c.TplName = "historic.html"
}

//WS 链接websocket
func (c *CreateController) WS() {
	beego.Info("进入WS1")

	
	// 检验http头中upgrader属性，若为websocket，则将http协议升级为websocket协议
	conn, err := (&websocket.Upgrader{}).Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	
	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error("Not a websocket connection")
		http.Error(c.Ctx.ResponseWriter, "Not a websocket hanfshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup Websocket connection:", err)
		return 
	}
	var version Version
	version.conn = conn
	services := models.GetAllService()

	data, err := json.Marshal(services)
	err = version.conn.WriteMessage(websocket.TextMessage ,data)
	beego.Info(services)

	defer func() {
		version.conn.Close()
	}()
	// 由于WebSocket一旦连接，便可以保持长时间通讯，则该接口函数可以一直运行下去，直到连接断开
	for {
		// 读取消息。如果连接断开，则会返回错误
		_, msgStr, err := version.conn.ReadMessage()
		if err != nil {
			break
		}
		var jsonversion  models.Version    //创建版本所需要的信息
		var jsonmessage []message         //接收解析之后得到表单数据
		var jsonservices = make([]models.SerList, 0)  //将服务列表的需要的服务名，和服务号封装起来，方便数据在前端和后端之间传输
		var jsonservice models.SerList   //服务列表的信息
		json.Unmarshal(msgStr, &jsonmessage)
		//beego.Info(jsonmessage)
		//手动的将一些信息赋值给版本对象的某些字段中
		for _, val := range jsonmessage {
			if val.Name == "versionnumber" {
				jsonversion.VersionNumber = val.Value
			} else if val.Name == "servicelist" {
				jsonservice.ServiceName, jsonservice.ServiceNumber = "", ""
				flag := 0
				for _, ch := range val.Value {
					if ch == '&' {
						flag = 1
						continue
					}
					if flag == 0 {
						jsonservice.ServiceName += string(ch)
					} else {
						jsonservice.ServiceNumber += string(ch)
					}
				}
				jsonservices = append(jsonservices, jsonservice)
			} else {
				jsonversion.Comment = val.Value
			}
		}
		jsonversion.ServiceList = jsonservices
		
		beego.Info("WS1-----------receive: ")
		beego.Info(jsonversion)
		create <- jsonversion
		

	}
}
// WS 修改的websocket的连接
func (c *ModifyController) WS() {

	beego.Info("进入WS2")
	// 检验http头中upgrader属性，若为websocket，则将http协议升级为websocket协议
	conn, err := (&websocket.Upgrader{}).Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	
	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error("Not a websocket connection")
		http.Error(c.Ctx.ResponseWriter, "Not a websocket hanfshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup Websocket connection:", err)
		return 
	}

	//websocket的连接
	var version Version
	version.conn = conn
	//定义一个封装 后台服务列表和后台版本信息的对象
	var sendmessage1 sendmessage

	//获取后台服务列表的信息
	services := models.GetAllService()
	
	//获取后台未发布的版本信息
	versions := models.GetAllStatus()

	sendmessage1.Services = services
	sendmessage1.Versions = versions
	data, err := json.Marshal(sendmessage1)
	checkErr(err)
	err = version.conn.WriteMessage(websocket.TextMessage, data)
	checkErr(err)

	beego.Info("服务端已经发送信息到客户端")
	defer version.conn.Close()

	for {
		_, msgStr, err := version.conn.ReadMessage()
		if err != nil {
			break
		}
		var jsonversion models.Version                 //修改的版本
		var jsonmessage []message                     //接收解析的json表单
		var jsonservices = make([]models.SerList, 0) //存储服务列表的信息，方便更新版本号后再修改到数据库
		var jsonservice models.SerList
		
		err = json.Unmarshal(msgStr, &jsonmessage)
		for _, val := range jsonmessage {
			if val.Name == "versionnumber" {
				id, _ := strconv.Atoi(val.Value)
				jsonversion.ID = id
			} else if val.Name == "servicelist" {
				//解析服务列表结构体中的数据，仅含有服务名和服务号
				jsonservice.ServiceName, jsonservice.ServiceNumber = "", ""
				flag := 0
				for _, ch := range val.Value {
					if ch == '&' {
						flag = 1
						continue
					}
					if flag == 0 {
						jsonservice.ServiceName += string(ch)
					} else {
						//处理更新时，绑定的小服务的版本小更新
						//if()
						jsonservice.ServiceNumber += string(ch)
					}
				}
				jsonservices = append(jsonservices, jsonservice)
			} else {
				jsonversion.Comment = val.Value
			}
		}
		jsonversion.ServiceList = jsonservices
		var ver = models.GetByID(jsonversion.ID)
		jsonversion.VersionNumber = ver.VersionNumber         
		beego.Info("WS2-----------receive: " + string(msgStr))
		update <- jsonversion
	}
}

// WS 发布websocket得连接
func (c *IssueController) WS() {
	beego.Info("进入WS3")
	conn, err := (&websocket.Upgrader{}).Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error("Not a websocket connection")
		http.Error(c.Ctx.ResponseWriter, "Not a websocket hanfshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup Websocket connection:", err)
		return
	}

	var version Version
	version.conn = conn

	//定义一个封装 后台服务列表和后台版本信息的对象
	var sendmessage1 sendmessage

	//获取后台服务列表的信息
	services := models.GetAllService()
	
	//获取后台未发布的版本信息
	versions := models.GetAllStatus()

	sendmessage1.Services = services
	sendmessage1.Versions = versions
	data, err := json.Marshal(sendmessage1)
	checkErr(err)
	err = version.conn.WriteMessage(websocket.TextMessage, data)
	checkErr(err)
	
	beego.Info("服务端已发送信息到客户端")
	defer version.conn.Close()

	for {
		_, msgStr, err := version.conn.ReadMessage()
		if err != nil {
			break
		}

		var jsonversion models.Version
		var jsonmessage []message
		err = json.Unmarshal(msgStr, &jsonmessage)
		id, _ := strconv.Atoi(jsonmessage[0].Value)
		jsonversion = models.GetByID(id)
		beego.Info(jsonversion.ID)
		beego.Info("WS3-----------receive: " + string(msgStr))
		issue <- jsonversion

	}
	
}
// WS 历史版本下websocket连接
func (c *HistoricController) WS() {
	beego.Info("进入WS4")
	conn, err := (&websocket.Upgrader{}).Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		beego.Error("Not a websocket connection")
		http.Error(c.Ctx.ResponseWriter, "Not a websocket hanfshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup Websocket connection:", err)
		return
	}

	var version Version
	version.conn = conn

	//定义一个封装 后台服务列表和后台版本信息的对象
	var sendmessage1 sendmessage

	services := models.GetAllService()
	versions := models.GetAllStatus1()

	sendmessage1.Services = services
	sendmessage1.Versions = versions
	
	data, err := json.Marshal(sendmessage1)
	checkErr(err)
	err = version.conn.WriteMessage(websocket.TextMessage, data)
	checkErr(err)

	beego.Info("服务端已发送信息到客户端")
	defer version.conn.Close()
	
}

