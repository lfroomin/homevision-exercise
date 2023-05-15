// This is a take home interview for HomeVision that focuses primarily on writing
// clean code that accomplishes a very practical task. We have a simple paginated
// API that returns a list of housesService along with some metadata. Your challenge is
// to write a script that meets the requirements.
//
// Note: this is a flaky API! That means that it will likely fail with a non-200
// response code. Your code must handle these errors correctly so that all photos
// are downloaded.
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type housesServiceInt interface {
	getPage(baseURL string, page, perPage int, housesCh chan<- []House) error
	savePhoto(store io.Writer, photoURL string) error
}

func main() {
	appCfg := loadConfig(".")

	var hs = housesService{
		httpClient: &http.Client{Timeout: time.Duration(10) * time.Second},
	}

	var doneCh = make(chan struct{})
	var errorCh = make(chan error)
	var houseCh = make(chan []House, appCfg.NumPages)

	// Any error will be placed on errorCh and will stop processing
	// TODO: This should be more robust for handling errors.
	go func() {
		e, ok := <-errorCh
		if ok {
			log.Printf("error occurred: %v\n", e)
			doneCh <- struct{}{}
		}
	}()

	go getHouses(hs, appCfg.HouseServiceUrl, appCfg.NumPages, appCfg.NumPerPage, houseCh, errorCh)

	go saveAllPhotos(hs, houseCh, doneCh, errorCh)

	// Wait for processing to complete
	<-doneCh

	close(errorCh)
	close(doneCh)
}

// getHouses invokes the House Service to get the data for each house. The retrieval
// utilizes paging, so the number of pages and number of houses per page are used to
// control the paging and ensure all houses are retrieved. The house data is placed
// onto the houseCh channel so that the data can be processed concurrently to save
// the house photo. Errors are placed on the errorCh channel.
func getHouses(hs housesServiceInt, housesURL string, numPages, numPerPage int, houseCh chan<- []House, errorCh chan<- error) {
	var wg = &sync.WaitGroup{}

	for i := 0; i < numPages; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := hs.getPage(housesURL, i+1, numPerPage, houseCh)
			if err != nil {
				errorCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(houseCh)
}

// saveAllPhotos retrieves the photo for each house and saves the image to a file. The
// house data is read from the houseCh channel which is a []House. So, the []House is
// iterated over to retrieve the house photo and save it to a file. The doneCh channel
// is signaled when all houses have been processed. Errors are placed on the errorCh
// channel.
func saveAllPhotos(hs housesServiceInt, houseCh <-chan []House, doneCh chan<- struct{}, errorCh chan<- error) {
	var wg = &sync.WaitGroup{}
	for houses := range houseCh {
		for _, h := range houses {
			wg.Add(1)
			go func(h House) {
				defer wg.Done()
				photoFile, err := os.Create(photoFilename(h))
				if err != nil {
					errorCh <- err
					return
				}

				defer func(photoFile *os.File) {
					err := photoFile.Close()
					if err != nil {
						errorCh <- err
						return
					}
				}(photoFile)

				err = hs.savePhoto(photoFile, h.PhotoURL)
				if err != nil {
					errorCh <- err
				}
			}(h)
		}
	}
	wg.Wait()
	doneCh <- struct{}{}
}
