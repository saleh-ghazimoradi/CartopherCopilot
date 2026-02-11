package cart

import "time"

type AddToCartResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Id        int       `json:"id"`
		UserId    int       `json:"user_id"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
	Error string `json:"error"`
}

type ViewCartResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Id        int `json:"id"`
		UserId    int `json:"user_id"`
		CartItems []struct {
			Id      int `json:"id"`
			Product struct {
				Id          int     `json:"id"`
				CategoryId  int     `json:"category_id"`
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Price       float64 `json:"price"`
				Stock       int     `json:"stock"`
				Sku         string  `json:"sku"`
				IsActive    bool    `json:"is_active"`
				Category    struct {
					Id          int       `json:"id"`
					Name        string    `json:"name"`
					Description string    `json:"description"`
					IsActive    bool      `json:"is_active"`
					CreatedAt   time.Time `json:"created_at"`
					UpdatedAt   time.Time `json:"updated_at"`
				} `json:"category"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"product"`
			Quantity  int       `json:"quantity"`
			Subtotal  float64   `json:"subtotal"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"cart_items"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"data"`
	Error string `json:"error"`
}
