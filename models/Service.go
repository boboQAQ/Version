package models

import (
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	
)
// Service 服务号的数据结构
type Service struct {
	ServiceID int 				`json:"serviceid"`
	ServiceName string 			`json:"servicename"`
	ServiceNumber string 		`json:"servicenumber"`
	NumberModifyTime time.Time 	`json:"Numbermodifytime"`
	ServiceStatus bool 			`json:"servicestatus"`

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