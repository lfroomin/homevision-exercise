package main

import (
	"fmt"
	"net/url"
	"strings"
)

// photoFilename returns the name of the file to store the house photo according
// to the format [id]-[address].[ext]
func photoFilename(house House) string {
	ext := getFileExtension(house.PhotoURL)
	// TODO: A check should be done to ensure the result is a valid file name
	// (ie doesn't include invalid characters, doesn't exceed length, etc.)
	fileName := fmt.Sprintf("%d-%s.%s", house.Id, house.Address, ext)
	return fileName
}

// getFileExtension returns the file extension from the house photo URL. If no
// file extension exists, 'jpg' is used by default.
func getFileExtension(photoURL string) string {
	// Default filename extension to 'jpg'
	ext := "jpg"
	u, err := url.Parse(photoURL)
	if err != nil {
		return ext
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return ext
	}
	return u.Path[pos+1 : len(u.Path)]
}
