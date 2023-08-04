package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/knhn1004/WhatDoIEat/internal/config"
	"github.com/knhn1004/WhatDoIEat/internal/models"
)

const yelpAPIEndpoint = "https://api.yelp.com/v3/businesses/search"

func GetRestaurants(restaurantType, location string, options map[string]string) ([]models.Restaurant, error) {
	if location == "" {
		location = "San Jose"
	}

	params := url.Values{}
	params.Add("term", restaurantType)
	params.Add("location", location)

	for key, value := range options {
		params.Add(key, value)
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", yelpAPIEndpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+config.YelpAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d %s. Response body: %s", resp.StatusCode, resp.Status, buf.String())
	}

	var yelpResp models.YelpResponse
	err = json.Unmarshal(buf.Bytes(), &yelpResp)
	if err != nil {
		return nil, err
	}

	return yelpResp.Businesses, nil
}
