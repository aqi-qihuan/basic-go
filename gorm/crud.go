package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Product 定义一个 Product 结构体，包含 gorm.Model 和 Code、Price 字段
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 打开调试模式
	db = db.Debug()

	// 迁移 schema
	// 初始化你的表结构
	db.AutoMigrate(&Product{})

	// Create
	// 创建一个 Product 实例，并插入到数据库中
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	// 定义一个 Product 实例
	var product Product
	// 根据整型主键查找
	db.First(&product, 1)
	// 查找 code 字段值为 D42 的记录
	db.First(&product, "code = ?", "D42")

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	// WHERE 条件是什么？
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}
