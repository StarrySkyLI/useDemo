package migrate

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"useDemo/application/rpc_demo/internal/dao/schema"
)

func Handle() {
	// 配置 MySQL 连接参数（替换为你的实际值）
	user := "root"
	password := "Newroot1515!"
	host := "localhost" // Kubernetes 服务名或 ClusterIP
	port := "3306"
	dbName := "demo" // 要自动创建的数据库名

	// 构造 DSN（注意此处连接的是默认数据库，如 `mysql`）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port)

	// 1. 连接到 MySQL 服务器（不指定数据库）
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接 MySQL 失败: " + err.Error())
	}

	// 2. 检查并创建数据库（如果不存在）
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)
	err = db.Exec(createDBQuery).Error
	if err != nil {
		panic("创建数据库失败: " + err.Error())
	}
	fmt.Printf("数据库 `%s` 已就绪\n", dbName)

	// 3. 重新连接到目标数据库
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)
	db, err = gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		panic("连接到目标数据库失败: " + err.Error())
	}

	// 4. 自动迁移表结构（示例）

	err = db.AutoMigrate(&schema.Game{})
	if err != nil {
		panic("迁移表失败: " + err.Error())
	}

	fmt.Println("数据库连接和初始化完成！")
}
