package panodata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// A Photo in the Panoramio database.
type Photo struct {
	Id         int64   `json:"photo_id"`
	Title      string  `json:"photo_title"`
	Url        string  `json:"photo_url"`
	FileUrl    string  `json:"photo_file_url"`
	Lat        float64 `json:"latitude"`
	Long       float64 `json:"longitude"`
	Width      int64   `json:"width"`
	Height     int64   `json:"height"`
	UploadDate string  `json:"upload_date"`
	OwnerId    int64   `json:"owner_id"`
	OwnerName  string  `json:"owner_name"`
	OwnerUrl   string  `json:"owner_url"`
}

// Run a Query returning a []Photo with the results. Retrieves all the photos'
// metadata by running nc concurrent queries to the web service, each
// requesting the maximum allowed amount of photos. Returns an error in the case
// of network or parsing errors.
func (query Query) Run(nc int) (photos []Photo, err error) {
	// Get the base URL
	url := query.URL()

	// Request the first 100 (or less) photos
	body, err := doRequest(addRange(url, 0, 100))
	if err != nil {
		return
	}

	// No additional request needed
	photos = body.Photos
	if int(body.Count) <= 100 && len(photos) <= 100 {
		return
	}

	type qRange struct {
		from, to int
	}

	// Run concurrent requests
	var wg sync.WaitGroup
	var mut sync.Mutex
	in := make(chan qRange, nc)

	// Start nConcurrent workers
	wg.Add(nc)
	for i := 0; i < nc; i++ {
		go func() {
			defer wg.Done()

			for r := range in {
				if body, err := doRequest(addRange(url, r.from, r.to)); err == nil {
					mut.Lock()
					photos = append(photos, body.Photos...)
					mut.Unlock()
				}
			}
		}()
	}

	// Send work
	go func() {
		for from := 100; from < int(body.Count); from += 100 {
			in <- qRange{from, from + 100}
		}
		close(in)
	}()

	// Wait results
	wg.Wait()

	return
}

// response defines the schema of a Panoramio API response.
type parsedResp struct {
	Count  int64   `json:"count"`
	Photos []Photo `json:"photos"`
}

func doRequest(url string) (*parsedResp, error) {
	var parsed parsedResp

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("panodata: http errror %v", err)
	}
	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c >= 300 {
		return nil, fmt.Errorf("panodata: request returned status code %v", c)
	}

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&parsed); err != nil {
		return nil, fmt.Errorf("panodata: json parsing error %v", err)
	}
	return &parsed, nil
}

func addRange(url string, from, to int) string {
	return fmt.Sprintf("%s&from=%d&to=%d", url, from, to)
}
