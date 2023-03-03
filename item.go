package fditch

import (
	"encoding/json"
)

// The Item struct represents the data for an itch.io item.
type Item struct {
	ID          string `json:"id"`
	Link        string `json:"link"`
	ImgLink     string `json:"img_link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	SalesLink   string `json:"sales_link"`
	EndDate     string `json:"end_date"`
	Genre       string `json:"genre"`
}

// ToJSON converts the Item to a JSON string.
// Returns an error if any.
func (item *Item) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(&item)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
