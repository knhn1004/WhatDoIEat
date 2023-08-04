package services

import (
	"net/url"

	"github.com/knhn1004/WhatDoIEat/internal/models"
)

func GenGoogleMapsURL(loc models.Location) string {
	baseURL := "https://www.google.com/maps/search/?api=1&query="

	// Concatenate address details
	address := loc.Address1 + " " + loc.Address2 + " " + loc.Address3 + " " + loc.City + " " + loc.ZipCode + " " + loc.State + " " + loc.Country

	// URL encode the address
	encodedAddress := url.QueryEscape(address)

	return baseURL + encodedAddress
}
