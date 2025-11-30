package db

import (
	"log"
)

// CreateTables 创建数据库表
func CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(200) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			image VARCHAR(500),
			stock INT DEFAULT 0,
			series VARCHAR(50) DEFAULT '拉布布',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS cart_items (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			product_id INT NOT NULL,
			quantity INT DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_user_id (user_id),
			INDEX idx_product_id (product_id),
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			total_price DECIMAL(10,2) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,

		`CREATE TABLE IF NOT EXISTS order_items (
			id INT AUTO_INCREMENT PRIMARY KEY,
			order_id INT NOT NULL,
			product_id INT NOT NULL,
			quantity INT NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_order_id (order_id),
			INDEX idx_product_id (product_id),
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("Error executing query: %v", err)
			return err
		}
	}

	log.Println("Database tables created successfully")
	return nil
}

// SeedProducts 初始化商品数据
func SeedProducts() error {
	products := []struct {
		name        string
		description string
		price       float64
		image       string
		stock       int
	}{
		{"拉布布 经典款", "经典拉布布盲盒，随机款式", 59.00, "https://via.placeholder.com/300x300?text=拉布布经典款", 100},
		{"拉布布 限定款", "限定版拉布布盲盒，稀有款式", 89.00, "https://via.placeholder.com/300x300?text=拉布布限定款", 50},
		{"拉布布 隐藏款", "隐藏款拉布布盲盒，超稀有", 199.00, "https://via.placeholder.com/300x300?text=拉布布隐藏款", 10},
		{"拉布布 套装", "拉布布系列套装，包含多个款式", 299.00, "https://via.placeholder.com/300x300?text=拉布布套装", 30},
		{"拉布布 特别版", "特别版拉布布盲盒", 129.00, "https://via.placeholder.com/300x300?text=拉布布特别版", 25},
	}

	for _, p := range products {
		query := `INSERT INTO products (name, description, price, image, stock, series) 
				  VALUES (?, ?, ?, ?, ?, '拉布布')
				  ON DUPLICATE KEY UPDATE name=name`
		_, err := DB.Exec(query, p.name, p.description, p.price, p.image, p.stock)
		if err != nil {
			log.Printf("Error seeding product %s: %v", p.name, err)
		}
	}

	log.Println("Products seeded successfully")
	return nil
}
