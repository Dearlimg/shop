package dao

import (
	"database/sql"
	"shop/global/db"
	"shop/model"
)

// GetCartItems 获取用户的购物车项
func GetCartItems(userID int) ([]model.CartItem, error) {
	rows, err := db.DB.Query(
		`SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
			p.id, p.name, p.description, p.price, p.image, p.stock, p.series, p.created_at, p.updated_at
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		var product model.Product
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt,
			&product.ID, &product.Name, &product.Description, &product.Price, &product.Image, &product.Stock, &product.Series, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		item.Product = product
		items = append(items, item)
	}
	return items, nil
}

// GetCartItemByID 根据ID获取购物车项
func GetCartItemByID(cartItemID, userID int) (*model.CartItem, error) {
	var item model.CartItem
	err := db.DB.QueryRow(
		"SELECT id, user_id, product_id, quantity FROM cart_items WHERE id = ? AND user_id = ?",
		cartItemID, userID,
	).Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetCartItemByUserAndProduct 获取用户和商品的购物车项
func GetCartItemByUserAndProduct(userID, productID int) (*model.CartItem, error) {
	var item model.CartItem
	err := db.DB.QueryRow(
		"SELECT id, user_id, product_id, quantity FROM cart_items WHERE user_id = ? AND product_id = ?",
		userID, productID,
	).Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// AddCartItem 添加购物车项
func AddCartItem(userID, productID, quantity int) error {
	_, err := db.DB.Exec(
		"INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)",
		userID, productID, quantity,
	)
	return err
}

// UpdateCartItemQuantity 更新购物车项数量
func UpdateCartItemQuantity(cartItemID, quantity int) error {
	_, err := db.DB.Exec(
		"UPDATE cart_items SET quantity = ? WHERE id = ?",
		quantity, cartItemID,
	)
	return err
}

// DeleteCartItem 删除购物车项
func DeleteCartItem(cartItemID int) error {
	_, err := db.DB.Exec("DELETE FROM cart_items WHERE id = ?", cartItemID)
	return err
}
