package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"

	"github.com/knhn1004/WhatDoIEat/internal/config"
)

// Supabase client
var supabaseClient *supabase.Client

// Initialize Supabase client
func InitSupabase() {
	supabaseClient = supabase.CreateClient(config.SupabaseURL, config.SupabaseKey)

	fmt.Printf("supabase: %v\n", supabaseClient) // TODO: remove this
}

func uploadImageFromUrl(img string) (string, error) {
	// download image data
	resp, err := http.Get(img)
	if err != nil {
		return "", fmt.Errorf("error downloading image: %v", err)
	}
	defer resp.Body.Close()

	var imgBuffer bytes.Buffer
	_, err = io.Copy(&imgBuffer, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading image data: %v", err)
	}

	// get file extension
	imgParts := strings.Split(img, ".")
	imgType := imgParts[len(imgParts)-1] // get last part, which should be the file extension

	// random filename
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), imgType)
	supabaseClient.Storage.From("recipe-images").Upload(filename, &imgBuffer)

	supabaseUrl := os.Getenv("SUPABASE_URL")
	url := fmt.Sprintf("%s/storage/v1/object/public/recipe-images/%s", supabaseUrl, filename)

	return url, nil
}
