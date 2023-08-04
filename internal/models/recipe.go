package models

type Ingredient struct {
	Ingredient string `json:"ingredient"`
	Quantity   string `json:"quantity"`
}

type Recipe struct {
	Meal             string       `json:"meal"`
	Name             string       `json:"name"`
	ShortDescription string       `json:"short_description"`
	UserId           string       `json:"user_id"`
	ImageURL         string       `json:"image_url"`
	Date             string       `json:"date"`
	Ingredients      []Ingredient `json:"ingredients"`
	Steps            []string     `json:"steps"`
}
