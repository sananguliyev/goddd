package entity

type Recipe struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	PrepareTime  int8   `json:"prepare_time"`
	Difficulty   int8   `json:"difficulty"`
	IsVegetarian *bool  `json:"is_vegetarian"`
}
