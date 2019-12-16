package controllers

import (
	"strconv"
	"Demo/models"
	"github.com/gorilla/websocket"
)
// Version 导出的客户端包的结构体
type Version struct {
	conn *websocket.Conn //用户websocket链接
	version models.Version//用户名称
}
//VersionNumberUpdate 版本号的更新计算函数
//type1类型为0代表是仅更新
//type1类型为1代表是发布更新
func VersionNumberUpdate(str string, type1 int) string {
	versionnumber1 := ""
	versionnumber2 := ""
	versionnumber3 := ""
	cnt := 0
	for _, val := range str {
		if val == '.' {
			cnt++
			continue
		}
		if cnt == 2 {
			versionnumber3 += string(val)
		} else if cnt == 1 {
			versionnumber2 += string(val)
		} else {
			versionnumber1 += string(val)
		}
	}
	if type1 == 1 {
		num, err:= strconv.Atoi(versionnumber3)
		checkErr(err)
		num++
		versionnumber3 = strconv.Itoa(num)
		checkErr(err)
		str = versionnumber1 + "." + versionnumber2 + "." + versionnumber3
	} else {
		num, err := strconv.Atoi(versionnumber2)
		checkErr(err)
		num++
		versionnumber2 = strconv.Itoa(num)
		checkErr(err)
		str = versionnumber1 + "." + versionnumber2 + ".0"
	}
	return str

}

//ServiceListUpdate 服务列表的版本号的更新计算函数
//type1类型为0代表是仅更新
//type1类型为1代表是发布更新
func ServiceListUpdate(serlist []models.SerList, type1 int) []models.SerList {
	i := 0
	for _, value := range serlist {
		str := value.ServiceNumber
		versionnumber1 := ""
		versionnumber2 := ""
		versionnumber3 := ""
		cnt := 0
		for _, val := range str {
			if val == '.' {
				cnt++
				continue
			}
			if cnt == 2 {
				versionnumber3 += string(val)
			} else if cnt == 1 {
				versionnumber2 += string(val)
			} else {
				versionnumber1 += string(val)
			}
		}
		if type1 == 0 {
			num, err:= strconv.Atoi(versionnumber3)
			checkErr(err)
			num++
			versionnumber3 = strconv.Itoa(num)
			checkErr(err)
			str = versionnumber1 + "." + versionnumber2 + "." + versionnumber3
		} else {
			num, err := strconv.Atoi(versionnumber2)
			checkErr(err)
			num++
			versionnumber2 = strconv.Itoa(num)
			checkErr(err)
			str = versionnumber1 + "." + versionnumber2 + ".0"
		}
		serlist[i].ServiceNumber = str
		i++
	}
	return serlist
}