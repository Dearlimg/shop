package dao

import (
	"errors"
	"fmt"
	"time"

	"shop/global/redis"
	"shop/model"

	redisv9 "github.com/redis/go-redis/v9"
)

const (
	// CartKeyPrefix 购物车键前缀
	CartKeyPrefix = "cart:user:"
	// CartExpireTime 购物车过期时间（30天）
	CartExpireTime = 30 * 24 * time.Hour
)

// getCartKey 获取购物车Redis键
func getCartKey(userID int) string {
	return fmt.Sprintf("%s%d", CartKeyPrefix, userID)
}

// getCartItemKey 获取购物车项Redis键
func getCartItemKey(userID, productID int) string {
	return fmt.Sprintf("%s%d:product:%d", CartKeyPrefix, userID, productID)
}

// GetCartItems 获取用户的购物车项（从Redis）
func GetCartItemsFromRedis(userID int) ([]model.CartItem, error) {
	// 获取所有购物车项的键
	pattern := fmt.Sprintf("%s%d:product:*", CartKeyPrefix, userID)
	keys, err := redis.Client.Keys(redis.GetContext(), pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("获取购物车键失败: %w", err)
	}

	if len(keys) == 0 {
		return []model.CartItem{}, nil
	}

	var items []model.CartItem
	for _, itemKey := range keys {
		// 获取商品ID和数量
		data, err := redis.Client.HGetAll(redis.GetContext(), itemKey).Result()
		if err != nil {
			if err == redisv9.Nil {
				continue
			}
			return nil, fmt.Errorf("获取购物车项失败: %w", err)
		}

		if len(data) == 0 {
			continue
		}

		var productID, quantity int
		if _, err := fmt.Sscanf(data["product_id"], "%d", &productID); err != nil {
			continue
		}
		if _, err := fmt.Sscanf(data["quantity"], "%d", &quantity); err != nil {
			continue
		}

		// 获取商品信息（从数据库）
		product, err := GetProductByID(fmt.Sprintf("%d", productID))
		if err != nil || product == nil {
			// 如果商品不存在，删除该购物车项
			redis.Client.Del(redis.GetContext(), itemKey)
			continue
		}

		item := model.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
			Product:   *product,
		}
		items = append(items, item)
	}

	return items, nil
}

// GetCartItemByUserAndProductFromRedis 从Redis获取用户和商品的购物车项
func GetCartItemByUserAndProductFromRedis(userID, productID int) (*model.CartItem, error) {
	itemKey := getCartItemKey(userID, productID)

	data, err := redis.Client.HGetAll(redis.GetContext(), itemKey).Result()
	if err != nil {
		if err == redisv9.Nil {
			return nil, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var quantity int
	if _, err := fmt.Sscanf(data["quantity"], "%d", &quantity); err != nil {
		return nil, err
	}

	return &model.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}, nil
}

// AddCartItemToRedis 添加购物车项到Redis
func AddCartItemToRedis(userID, productID, quantity int) error {
	itemKey := getCartItemKey(userID, productID)
	cartKey := getCartKey(userID)

	// 使用Hash存储购物车项
	itemData := map[string]interface{}{
		"product_id": productID,
		"quantity":   quantity,
		"user_id":    userID,
	}

	// 设置购物车项
	err := redis.Client.HSet(redis.GetContext(), itemKey, itemData).Err()
	if err != nil {
		return fmt.Errorf("添加购物车项失败: %w", err)
	}

	// 设置过期时间
	redis.Client.Expire(redis.GetContext(), itemKey, CartExpireTime)
	redis.Client.Expire(redis.GetContext(), cartKey, CartExpireTime)

	return nil
}

// UpdateCartItemQuantityInRedis 更新Redis中购物车项数量
func UpdateCartItemQuantityInRedis(userID, productID, quantity int) error {
	itemKey := getCartItemKey(userID, productID)

	// 检查项是否存在
	exists, err := redis.Client.Exists(redis.GetContext(), itemKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("购物车项不存在")
	}

	// 更新数量
	err = redis.Client.HSet(redis.GetContext(), itemKey, "quantity", quantity).Err()
	if err != nil {
		return fmt.Errorf("更新购物车项失败: %w", err)
	}

	// 更新过期时间
	redis.Client.Expire(redis.GetContext(), itemKey, CartExpireTime)

	return nil
}

// DeleteCartItemFromRedis 从Redis删除购物车项
func DeleteCartItemFromRedis(userID, productID int) error {
	itemKey := getCartItemKey(userID, productID)
	return redis.Client.Del(redis.GetContext(), itemKey).Err()
}

// GetCartItemByIDFromRedis 根据ID获取购物车项（需要从productID反推）
// 注意：Redis版本中，我们使用productID作为标识，而不是自增ID
func GetCartItemByIDFromRedis(userID, productID int) (*model.CartItem, error) {
	return GetCartItemByUserAndProductFromRedis(userID, productID)
}

// ClearCartFromRedis 清空用户的购物车
func ClearCartFromRedis(userID int) error {
	pattern := fmt.Sprintf("%s%d:product:*", CartKeyPrefix, userID)
	keys, err := redis.Client.Keys(redis.GetContext(), pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return redis.Client.Del(redis.GetContext(), keys...).Err()
	}

	return nil
}

// GetCartItemWithProductFromRedis 从Redis获取购物车项及其商品信息
func GetCartItemWithProductFromRedis(userID, productID int) (*model.CartItem, *model.Product, error) {
	item, err := GetCartItemByUserAndProductFromRedis(userID, productID)
	if err != nil {
		return nil, nil, err
	}
	if item == nil {
		return nil, nil, nil
	}

	product, err := GetProductByID(fmt.Sprintf("%d", productID))
	if err != nil {
		return nil, nil, err
	}
	if product == nil {
		// 商品不存在，删除购物车项
		DeleteCartItemFromRedis(userID, productID)
		return nil, nil, nil
	}

	return item, product, nil
}
