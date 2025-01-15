package models

type ItemRequest struct {
	Name  string  `json:"name" example:"Sample Item"`
	Price float64 `json:"price" example:"19.99"`
}

type ItemResponse struct {
	ID    uint64  `json:"id" example:"1"`
	Name  string  `json:"name" example:"Sample Item"`
	Price float64 `json:"price" example:"19.99"`
}
