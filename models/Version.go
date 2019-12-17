package models

import (
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	
)

// SerList 服务列表的结构体
//对应服务和版本号
type SerList struct {
	ServiceName string 		`json:"servicename"`
	ServiceNumber string	`json:"servicenumber"`
}

//Version 版本结构体
type Version struct {
	ID int               	`json:"id"`
	VersionNumber string 	`json:"versionnumber"`
	ServiceList []SerList	`json:"servicelist"`	
	Status bool			 	`json:"status"`	
	IssueTime time.Time	 	`json:"issuetime"`
	CreatTime time.Time	 	`json:"creattime"`
	Comment string		 	`json:"comment"`
	IssueStatus int         `json:"issuestatus"`
}
func checkErr (err error) {
	if err != nil {
		panic(err)
	}
}

var (
	dbhostsip = "127.0.0.1:3306"//mysql登入点
	dbusername = "root"//mysql登入名
	dbpassword = "huangbo0.0"//mysql登入密码
	dbname     = "test"//链接的数据库名
	//depart是表名
)


// GetAllStatus 获取所有状态为未发布的版本
func GetAllStatus() []Version{
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&parseTime=true")
	checkErr(err)

	rows, err := db.Query("select * from versions where Status != true")
	checkErr(err)
	var versions = make([]Version, 0)
	for rows.Next() {
		var version Version
		var serlist SerList
		var serlists = make([]SerList,0)
		var str string
		flag := 0
		err = rows.Scan(&version.ID, &version.VersionNumber, &str, &version.Status, &version.IssueTime, &version.CreatTime, &version.Comment, &version.IssueStatus)
		for _, ch :=  range str {
			if ch == '<' || ch == '>' {
				if ch == '>' {
					serlists = append(serlists, serlist)
				} else {
					serlist.ServiceName, serlist.ServiceNumber = "", ""
					flag = 0
				}
				continue
			}
			if ch == ' ' {
				flag = 1
				continue
			}
			if flag ==0 {
				serlist.ServiceName += string(ch)
			} else {
				serlist.ServiceNumber += string(ch)
			}
		}
		version.ServiceList = serlists
		checkErr(err)
		versions = append(versions, version)
	}
	db.Close()
	return versions

}

// GetAllStatus1 获取所有状态为已经发布的历史版本
func GetAllStatus1() []Version{
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&parseTime=true")
	checkErr(err)

	rows, err := db.Query("select * from versions where Status = true")
	checkErr(err)
	var versions = make([]Version, 0)
	for rows.Next() {
		var version Version
		var serlist SerList
		var serlists = make([]SerList,0)
		var str string
		flag := 0
		err = rows.Scan(&version.ID, &version.VersionNumber, &str, &version.Status, &version.IssueTime, &version.CreatTime, &version.Comment, &version.IssueStatus)
		for _, ch :=  range str {
			if ch == '<' || ch == '>' {
				if ch == '>' {
					serlists = append(serlists, serlist)
				} else {
					serlist.ServiceName, serlist.ServiceNumber = "", ""
					flag = 0
				}
				continue
			}
			if ch == ' ' {
				flag = 1
				continue
			}
			if flag ==0 {
				serlist.ServiceName += string(ch)
			} else {
				serlist.ServiceNumber += string(ch)
			}
		}
		version.ServiceList = serlists
		checkErr(err)
		versions = append(versions, version)
	}
	db.Close()
	return versions

}

// GetByID 根据ID获取特定版本信息
func GetByID(ID int) Version {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&parseTime=true")
	checkErr(err)

	var version Version
	row := db.QueryRow("select * from versions where ID = ?", ID)
	var serlist SerList
	var serlists = make([]SerList,0)
	var str string
	flag := 0
	err = row.Scan(&version.ID, &version.VersionNumber, &str, &version.Status, &version.IssueTime, &version.CreatTime, &version.Comment, &version.IssueStatus)
	for _, ch :=  range str {
		if ch == '<' || ch == '>' {
			if ch == '>' {
				serlists = append(serlists, serlist)
			} else {
				serlist.ServiceName, serlist.ServiceNumber = "", ""
				flag = 0
			}
			continue
		}
		if ch == ' ' {
			flag = 1
			continue
		}
		if flag ==0 {
			serlist.ServiceName += string(ch)
		} else {
			serlist.ServiceNumber += string(ch)
		}
	}
	version.ServiceList = serlists
	checkErr(err)

	db.Close()
	return version
}
// CreatVersion 创建版本
func CreatVersion(version *Version) int64 {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&loc=Local")
	checkErr(err)

	stmt, err := db.Prepare("insert versions set  VersionNumber=?, ServiceList=?, Status=?, IssueTime=?, CreatTime=?, Comment=?, IssueStatus=?")
	checkErr(err)
	var servicelist string
	for _, val := range version.ServiceList {
		servicelist = servicelist + "<" + val.ServiceName + " " + val.ServiceNumber + ">"
	}

	res, err := stmt.Exec(version.VersionNumber, servicelist, version.Status, version.IssueTime, version.CreatTime, version.Comment, version.IssueStatus)
	checkErr(err)

	id, err := res.LastInsertId()
	db.Close()
	return id
}
// UpdateVersion 更新状态未发布的版本
func UpdateVersion(version *Version) {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&loc=Local")
	checkErr(err)

	stmt, err := db.Prepare("update versions set ServiceList=?, Comment=? where ID=?")
	checkErr(err)
	var servicelist string
	for _, val := range version.ServiceList {
		servicelist = servicelist + "<" + val.ServiceName + " " + val.ServiceNumber + ">"
	}

	_, err = stmt.Exec(servicelist, version.Comment, version.ID)

	checkErr(err)

	db.Close()
}

// IssueVersion 发布版本
func IssueVersion(version *Version) {
	db, err := sql.Open("mysql", dbusername+":"+dbpassword+"@tcp("+dbhostsip+")/"+dbname+"?charset=utf8&loc=Local")
	checkErr(err)

	stmt, err := db.Prepare("update versions set ServiceList=?, Status=?, IssueTime=? where ID=?")
	checkErr(err)
	var servicelist string
	for _, val := range version.ServiceList {
		servicelist = servicelist + "<" + val.ServiceName + " " + val.ServiceNumber + ">"
	}
	
	_, err = stmt.Exec(servicelist, version.Status, version.IssueTime, version.ID)
	checkErr(err)

	db.Close()
}

func init(){
	
}
	
