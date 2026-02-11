package orders

type OrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Id     int     `json:"id"`
		Status string  `json:"status"`
		Total  float64 `json:"total_amount"`
	} `json:"data"`
	Error string `json:"error"`
}
