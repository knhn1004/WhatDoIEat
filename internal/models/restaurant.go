package models

type Restaurant struct {
	Location Location `json:"location"`
	Name     string   `json:"name"`
	State    string   `json:"state"`
	Phone    string   `json:"phone"`
	URL      string   `json:"url"`
	ImageURL string   `json:"image_url"`
	Rating   float64  `json:"rating"`
}

type Location struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
	State    string `json:"state"`
	Country  string `json:"country"`
}
