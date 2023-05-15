package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_GetHouses(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		numPages int
		error    error
	}{
		{
			name:     "happy path",
			numPages: 2,
		},
		{
			name:     "error",
			numPages: 1,
			error:    errors.New("an error occurred"),
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := &housesIntStub{err: tc.error}

			housesCh := make(chan []House)
			errorCh := make(chan error)
			doneCh := make(chan struct{})

			var resErrors []error
			go func() {
				for {
					select {
					case e := <-errorCh:
						resErrors = append(resErrors, e)
					case <-doneCh:
						return
					}
				}
			}()

			getHouses(h, "", tc.numPages, 5, housesCh, errorCh)

			doneCh <- struct{}{}
			close(errorCh)

			assert.Equal(t, tc.numPages, h.getPageCalled)
			if tc.error == nil {
				assert.Nil(t, resErrors)
			} else {
				assert.Equal(t, []error{tc.error}, resErrors)
			}
		})
	}
}

func Test_SaveAllPhotos(t *testing.T) {
	// TODO: Implement test cases
}

type housesIntStub struct {
	err           error
	getPageCalled int
}

func (s *housesIntStub) getPage(_ string, _, _ int, _ chan<- []House) error {
	s.getPageCalled++

	if s.err != nil {
		return s.err
	}
	return nil
}

func (s *housesIntStub) savePhoto(_ io.Writer, _ string) error {
	if s.err != nil {
		return s.err
	}
	return nil
}
