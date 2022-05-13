package dtos

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
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
