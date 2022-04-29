package models

type Metadata struct {
	BaseModel
	Name  string `json:"name"`
	Value string `json:"value"`
}
