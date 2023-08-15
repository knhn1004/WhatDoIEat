package services

import (
	"fmt"

	"github.com/minoplhy/duckduckgo-images-api"
)

func GetImageUrlFromText(keyword string) (string, error) {
	var img string
	serp := goduckgo.Search(goduckgo.Query{Keyword: keyword})
	if len(serp.Results) > 0 {
		img = serp.Results[0].Image
	}
	if img == "" {
		return "", fmt.Errorf("no images found for recipe: %s", keyword)
	}

	return img, nil
}
