package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PhotoFilename(t *testing.T) {
	t.Parallel()

	house := House{
		Id:       123,
		Address:  "456 Main Street Any City, CA 91234",
		PhotoURL: "https://photos.com/home.jpg",
	}
	exp := "123-456 Main Street Any City, CA 91234.jpg"
	res := photoFilename(house)
	assert.Equal(t, exp, res)
}

func Test_GetFileExtension(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		photoURL string
		exp      string
	}{
		{
			name:     "happy path",
			photoURL: "/test-photo.jpg",
			exp:      "jpg",
		},
		{
			name:     "default jpg",
			photoURL: "/test-photo",
			exp:      "jpg",
		},
		{
			name:     "fail URL parse",
			photoURL: "1:",
			exp:      "jpg",
		},
		{
			name:     "non jpg",
			photoURL: "/test-photo.tiff",
			exp:      "tiff",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			res := getFileExtension(tc.photoURL)

			assert.Equal(t, tc.exp, res)
		})
	}
}
