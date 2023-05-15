package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// House is the structure of house data returned from the House Service.
type House struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

type housesService struct {
	httpClient *http.Client
}

// getPage invokes the House Service to retrieve a single page of house data. If the
// House Service returns an error, then the retrieval is retried (up to 10 times).
func (s housesService) getPage(housesURL string, page, perPage int, housesCh chan<- []House) error {
	url := fmt.Sprintf("%s?page=%d&per_page=%d", housesURL, page, perPage)
	log.Println("get house data page: ", url)

	var resp = &http.Response{}
	var err error

	// TODO: This approach is simplistic as it retries failed attempts to get house data a fixed number of
	// times with a fixed delay between attempts. In a production environment, this would be changed to be
	// more robust according to the needs defined for the feature.
	for tries, status := 0, http.StatusInternalServerError; tries < 10 && status != http.StatusOK; tries++ {
		resp, err = s.httpClient.Get(url)
		if err != nil {
			return fmt.Errorf("getPage: error HTTP get: %w", err)
		}

		status = resp.StatusCode
		if status != http.StatusOK {
			// Sleep before trying again
			time.Sleep(500 * time.Millisecond)
			log.Printf("get house data page unsuccessful status (%d); trying again", status)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("getPage: non-200 return status code: %d", resp.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("getPage: failed to close response body: %v\n", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("getPage: error reading response body: %w", err)
	}

	// The data is returned in a list with key 'houses', so this struct is
	// used to unmarshal the data.
	type HousesJson struct {
		Houses []House `json:"houses"`
	}

	var houses HousesJson
	err = json.Unmarshal(respBody, &houses)
	if err != nil {
		return fmt.Errorf("getPage: error unmarshaling response body: %w", err)
	}

	housesCh <- houses.Houses
	return nil
}

// savePhoto retrieves the house photo using the photoURL attributes of the house
// data. The photo image is stored to a file with a file name using the format
// [id]-[address].[ext].
func (s housesService) savePhoto(store io.Writer, photoURL string) error {
	log.Println("save photo: ", photoURL)

	resp, err := s.httpClient.Get(photoURL)
	if err != nil {
		return fmt.Errorf("savePhoto: error HTTP get: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("savePhoto: non-200 return status code: %d", resp.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("savePhoto: failed to close response body: %v\n", err)
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("savePhoto: error reading response body: %w", err)
	}

	_, err = store.Write(respBody)
	if err != nil {
		return fmt.Errorf("savePhoto: error writing photo to file: %w", err)
	}

	return nil
}
