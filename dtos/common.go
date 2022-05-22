package dtos

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ChangeStatusRequest struct {
	Status     string `json:"status"`
	CustomerID string `json:"customer_id"`
}

type VariantFilter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProductFilter struct {
	Variants  []VariantFilter `json:"variants,omitempty"`
	MinPrice  float64         `json:"min_price"`
	MaxPrice  float64         `json:"max_price"`
	MinRating float64         `json:"min_rating"`
}

type CustomerFilter struct {
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
	Name   string `json:"name"`
}
