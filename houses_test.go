package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_GetPage(t *testing.T) {
	testCases := []struct {
		name   string
		page   int
		exp    []House
		errMsg string
	}{
		{
			name: "happy path",
			page: 1,
			exp: []House{
				{Id: 1, Address: "123 Street Any City", PhotoURL: "https://photo-1.jpg"},
				{Id: 2, Address: "456 Street Any City", PhotoURL: "https://photo-2.jpg"},
			},
		},
		{
			name:   "error",
			page:   0,
			errMsg: "getPage: non-200 return status code: 500",
		},
		// TODO: Add tests for response body cases
	}

	svr := httptest.NewServer(http.HandlerFunc(testHandler))
	defer svr.Close()

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			houses := housesService{
				httpClient: &http.Client{Timeout: time.Duration(10) * time.Second},
			}

			housesCh := make(chan []House)
			done := make(chan struct{})
			var resHouses []House

			go func() {
				resHouses = <-housesCh
				done <- struct{}{}
			}()

			err := houses.getPage(svr.URL, tc.page, 5, housesCh)
			close(housesCh)

			<-done

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.exp, resHouses)
			}
		})
	}
}

func Test_SavePhoto(t *testing.T) {
	testCases := []struct {
		name     string
		photoURL string
		errMsg   string
	}{
		{
			name:     "happy path",
			photoURL: "/test-photo.jpg",
		},
		{
			name:     "error",
			photoURL: "/error",
			errMsg:   "savePhoto: non-200 return status code: 500",
		},
		// TODO: Add tests for response body cases
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Could be enhanced to handle additional test scenarios
		if r.URL.Path == "/error" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer svr.Close()

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			houses := housesService{
				httpClient: &http.Client{Timeout: time.Duration(10) * time.Second},
			}

			err := houses.savePhoto(io.Discard, svr.URL+tc.photoURL)

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Could be enhanced to handle additional test scenarios
	if r.URL.Query().Get("page") == "0" {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		data := struct {
			Houses []House `json:"houses"`
		}{
			Houses: []House{
				{
					Id:       1,
					Address:  "123 Street Any City",
					PhotoURL: "https://photo-1.jpg",
				},
				{
					Id:       2,
					Address:  "456 Street Any City",
					PhotoURL: "https://photo-2.jpg",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(data)
	}
}
