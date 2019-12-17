package models

import (
	"time"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	
)

// Service 传递给前端的服务列表的的数据结构
type Service struct {
	ServiceID int 				`json:"serviceid"`//项目ID
	ServiceName string 			`json:"servicename"`//项目Name
	ServiceNumber string 		`json:"servicenumber"`//标签Name
	NumberModifyTime time.Time 	`json:"Numbermodifytime"`//未定义
	ServiceStatus bool 			`json:"servicestatus"`//未定义

}
// Projects 接收调用gitlab的接口 获取一个群组下的所有服务工程的结构体
type Projects struct {
	ID int `json:"id"`
	Name string `json:"name"`
}


// Tags 接收调用gitlab的接口 获取一个project下的所有标签的结构体
type Tags struct {
	Commits Commit `json:"commit"`
	Name string `json:"name"`  //标签名
	Target string `json:"target"`
	Message string `json:"message"` //该标签的一个解释短句
}
//Commit 接收commit信息的结构体，
type Commit struct {
	ID string `json:"id"` 			//用于合并之后打标签使用的ref 
    ShortID string `json:"short_id"`//使用提交后的由安全散列算法得到的一个固定长度的字符串
    Title string `json:"title"`
    CreatedAt time.Time `json:"created_at"`
    ParentIds string `json:"parent_ids"`
    Message string `json:"message"`
    AuthorName string `json:"author_name"`
    AuthorEmail string `json:"author_email"`
    CommitterName string `json:"committer_name"`
    CommitterEmail string `json:"committer_email"`
    CommittedDate time.Time `json:"committed_date"`
}

//HTTPGetProjects 调用get接口获取812群组下的所有服务id及名字
func HTTPGetProjects() []Projects {

    client := &http.Client{}
	url := "https://git.ucloudadmin.com/api/v4/groups/812/projects"
	//新建一个请求，参数为方法，路径，body参数
	res, err := http.NewRequest("GET",url, nil)
	//添加一个 --header
    res.Header.Add("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
    response, _ := client.Do(res)

	checkErr(err)
	defer response.Body.Close()
	//读出返回的响应的主体
	body, err := ioutil.ReadAll(response.Body)
	checkErr(err)

    var data []Projects
	json.Unmarshal(body, &data)

	beego.Info("json解析后得到"  )
	beego.Info( data)
	
	return data
}
//HTTPGetProject 调用get接口获取项目id为id的项目信息
func HTTPGetProject(id int) Projects {

    client := &http.Client{}
	url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id)
	//新建一个请求，参数为方法，路径，body参数
	res, err := http.NewRequest("GET",url, nil)
	//添加一个 --header
    res.Header.Add("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
    response, _ := client.Do(res)

	checkErr(err)
	defer response.Body.Close()
	//读出返回的响应的主体
	body, err := ioutil.ReadAll(response.Body)
	beego.Info(string(body))
	checkErr(err)

    var data Projects
	json.Unmarshal(body, &data)

	beego.Info("json解析后得到"  )
	beego.Info( data)
	
	return data
}
// GetCommitSha 获取某个项目的最新的sha
func GetCommitSha(id int) string {
	client := &http.Client{}
	
	url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id) + "/repository/commits"

	res, err := http.NewRequest("GET", url, nil)
	checkErr(err)
	res.Header.Add("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")

	resp, err := client.Do(res)
	checkErr(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	var commits []Commit
	json.Unmarshal(body, &commits)
	return commits[0].ID

}

//HTTPGetTags 调用get接口获取id项目组下的所有标签的信息，按名称以相反的字母顺序排序，并返回第一个
func HTTPGetTags(id int) Tags {

    client := &http.Client{}
    url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id) + "/repository/tags"

    res, err := http.NewRequest("GET",url, nil)
    res.Header.Add("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
    response, _ := client.Do(res)

	checkErr(err)
    defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	checkErr(err)


	var data []Tags
	var data1 Tags
	json.Unmarshal(body, &data)
	if len(data) == 0 {
		_, data1 = HTTPPostTag(id, "v0.0.1")
	} else {
		data1 = data[0]
	}
	
	return data1
}
//HTTPPostTag 创建标签
func HTTPPostTag(id int, tagname string) (string, Tags) {
	client := &http.Client{}
	url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id) + "/repository/tags"
	bd := "&tag_name="+tagname+"&ref=master&message=调用API后，给master打上标签"
    req, err := http.NewRequest("POST", url, strings.NewReader(bd))
    checkErr(err)
 
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
 
	resp, err := client.Do(req)
	checkErr(err)
    defer resp.Body.Close()
 
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	var errmessage ErrMessage
	var data1 Tags
	json.Unmarshal(body, &data1)
	json.Unmarshal(body, &errmessage)

	return errmessage.Message, data1

}






//GetAllService 获取全部未禁用的服务
func GetAllService() []Service {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&parseTime=true")
	checkErr(err)

	rows, err := db.Query("select * from services where ServiceStatus != true")
	checkErr(err)

	var services = make([]Service, 0)
	var service Service
	for rows.Next() {
		err = rows.Scan(&service.ServiceID, &service.ServiceName, &service.ServiceNumber, &service.NumberModifyTime, &service.ServiceStatus)
		services = append(services, service)
	}
	db.Close()
	return services

}

//UpdateService 更新服务的状态是否禁用或服务版本号更新
func UpdateService(flag int, service *Service) {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&loc=Local")
	checkErr(err)
	if flag == 0 {
		stmt, err := db.Prepare("update services set ServiceNumber=? where ID=?")
		checkErr(err)

		_, err = stmt.Exec(service.ServiceNumber, service.ServiceID)
		checkErr(err)
	} else {
		stmt, err := db.Prepare("update services set ServiceStatus=? where ID=?")
		checkErr(err)

		_, err = stmt.Exec(service.ServiceStatus, service.ServiceID)
		checkErr(err)
	}

	db.Close()
}

//InsertService 插入新服务
func InsertService(service *Service) {
	//db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&loc=Local")
	//checkErr(err)

	//stmt, err := db.Prepare("insert services set ServiceName=?, ServiceNumber=?, NumberModifyTime=?, ServiceStatus=? ")
}