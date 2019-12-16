package models

import (
	"time"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/astaxie/beego"
	
)
//MergeRequest 接收与合并请求有关的结构体
type MergeRequest struct {
	ID int `json:"id"`
	IID int `json:"iid"`
	ProjectID int `json:"project_id"`
	Title string `json:"title"`
	MergedAt time.Time `json:"merged_at"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
}

// ErrMessage 接收报错信息
type ErrMessage struct {
	Message string `json:"message"`
}

//ListMergeRequest 列出项目下所有的合并请求
func ListMergeRequest(id int) []MergeRequest {
	client := &http.Client{}
	url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id) + "/merge_requests"
	//创建一个新的请求
    req, err := http.NewRequest("GET", url, nil)
	checkErr(err)

	//给调用的url添加Header
    req.Header.Set("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
 
	//使用client.Do执行请求，获得一个respond响应
	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	//读取返回的响应中的信息
    body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	beego.Info(body)
	var mergerequest []MergeRequest
	json.Unmarshal(body, &mergerequest)
	return mergerequest
}

//CreateMergeRequest 创造合并请求函数
func CreateMergeRequest(id int) MergeRequest{
	client := &http.Client{}
	url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id) + "/merge_requests"
	bd  :="id=" + strconv.Itoa(id) + "&source_branch=develop&target_branch=master&title=创建合并请求"
	//创建一个新的合并请求
    req, err := http.NewRequest("POST", url, strings.NewReader(bd))
	checkErr(err)

	//给调用的url添加Header
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
 
	//使用client.Do执行请求，获得一个respond响应
	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	//读取返回的响应中的信息
    body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	beego.Info(body)
	var mergerequest MergeRequest
	json.Unmarshal(body, &mergerequest)
	return mergerequest
}
// Merge 确认合并
func Merge(id, iid int) string {
	client := &http.Client{}

	url := "https://git.ucloudadmin.com/api/v4/projects/" + strconv.Itoa(id) + "/merge_requests/" + strconv.Itoa(iid) + "/merge"
    req, err := http.NewRequest("PUT", url, nil)
	checkErr(err)
 
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("PRIVATE-TOKEN", "vs68e5weD4Z9gSwWEA8u")
 
	resp, err := client.Do(req)
	checkErr(err)
    defer resp.Body.Close()
 
    body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	var errmessage ErrMessage
	json.Unmarshal(body, &errmessage)

	return errmessage.Message

}