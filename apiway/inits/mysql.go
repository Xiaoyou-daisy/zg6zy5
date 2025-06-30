package inits

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error
var DB *gorm.DB

func InitMysql() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "zg6project:2003225zyh@tcp(14.103.243.149:3306)/zg6project?charset=utf8mb4&parseTime=True&loc=Local"
	// 使用 GORM 打开数据库连接
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")

	} else {
		// 迁移 schema
		fmt.Println("connect database success")
	}
}
