package db

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"reflect"
	"time"
)

// DbEngine 数据库对象
var DbEngine *gorm.DB

// CreateTableIfNotExist 创建表
func CreateTableIfNotExist(Engine *gorm.DB, tableModels []interface{}) {
	for _, value := range tableModels {
		err := Engine.AutoMigrate(value)
		if err != nil {
			fmt.Println("Create table ", reflect.TypeOf(value), " error!")
		}
	}
}

// InitDB 初始化数据库
func InitDB() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(10*time.Second))
	go func(ctx context.Context) {
		// 从环境变量中获取数据库配置
		Username := os.Getenv("ACV_DB_USERNAME")
		Password := os.Getenv("ACV_DB_PASSWORD")
		Hostname := os.Getenv("ACV_DB_HOSTNAME")
		Dbname := os.Getenv("ACV_DB_DBNAME")
		// 连接数据库
		connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4,utf8",
			Username, Password, Hostname, Dbname)
		var err1 error
		DbEngine, err1 = gorm.Open(mysql.Open(connStr), &gorm.Config{})
		if err1 != nil {
			panic(any("Database connect error," + err1.Error()))
		}
		sqlDB, err := DbEngine.DB()
		if err != nil {
			panic(any("Database error"))
		}
		var temp []interface{}
		temp = append(temp)
		CreateTableIfNotExist(DbEngine, temp)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(10000)
		sqlDB.SetConnMaxLifetime(time.Second * 3)
		cancel()
	}(ctx)
	fmt.Println(ctx)
	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			fmt.Println("context timeout exceeded")
			panic(any("Timeout when initialize database connection"))
		case context.Canceled:
			fmt.Println("context cancelled by force. wblog process is complete")
		}
	}
}
