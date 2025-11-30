package db

import (
	"log"

	"shop/model"
)

// CreateTables 创建数据库表
func CreateTables() error {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
	)
	if err != nil {
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

// SeedProducts 初始化商品数据
func SeedProducts() error {
	products := []model.Product{
		{Name: "拉布布 经典款", Description: "经典拉布布盲盒，随机款式", Price: 59.00, Image: "https://via.placeholder.com/300x300?text=拉布布经典款", Stock: 100, Series: "拉布布"},
		{Name: "拉布布 限定款", Description: "限定版拉布布盲盒，稀有款式", Price: 89.00, Image: "https://via.placeholder.com/300x300?text=拉布布限定款", Stock: 50, Series: "拉布布"},
		{Name: "拉布布 隐藏款", Description: "隐藏款拉布布盲盒，超稀有", Price: 199.00, Image: "https://via.placeholder.com/300x300?text=拉布布隐藏款", Stock: 10, Series: "拉布布"},
		{Name: "拉布布 套装", Description: "拉布布系列套装，包含多个款式", Price: 299.00, Image: "https://via.placeholder.com/300x300?text=拉布布套装", Stock: 30, Series: "拉布布"},
		{Name: "拉布布 特别版", Description: "特别版拉布布盲盒", Price: 129.00, Image: "https://via.placeholder.com/300x300?text=拉布布特别版", Stock: 25, Series: "拉布布"},
	}

	for _, p := range products {
		// 使用 FirstOrCreate 避免重复插入
		var existingProduct model.Product
		result := DB.Where("name = ?", p.Name).First(&existingProduct)
		if result.Error != nil {
			// 如果不存在则创建
			if err := DB.Create(&p).Error; err != nil {
				log.Printf("Error seeding product %s: %v", p.Name, err)
			}
		}
	}

	log.Println("Products seeded successfully")
	return nil
}
