package dao

import (
	"database/sql"
	"shop/global/db"
	"shop/model"
)

// GetProducts 获取所有商品
func GetProducts() ([]model.Product, error) {
	rows, err := db.DB.Query("SELECT id, name, description, price, image, stock, series, created_at, updated_at FROM products ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Image, &p.Stock, &p.Series, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID 根据ID获取商品
func GetProductByID(id string) (*model.Product, error) {
	var product model.Product
	err := db.DB.QueryRow(
		"SELECT id, name, description, price, image, stock, series, created_at, updated_at FROM products WHERE id = ?",
		id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetProductStock 获取商品库存
func GetProductStock(productID int) (int, error) {
	var stock int
	err := db.DB.QueryRow("SELECT stock FROM products WHERE id = ?", productID).Scan(&stock)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return stock, err
}

// UpdateProductStock 更新商品库存
func UpdateProductStock(productID, quantity int) error {
	_, err := db.DB.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", quantity, productID)
	return err
}
