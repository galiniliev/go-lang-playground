package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"golang.org/x/net/http2"
)

const TrackingId = "x-ms-tracking-id"

var requestTrace map[string]resty.TraceInfo

func TestParallel() {
	var wg sync.WaitGroup

	// Create a Resty Client
	client := GetRestyClient()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.R().
				EnableTrace().
				Get("https://apim-mckq3zpiubjie.azure-api.net/mock/get")

			if err != nil {
				fmt.Println("Error:", err)
			} else {
				// defer resp.Body.Close()
				fmt.Printf("Status Code: %v , Proto: %v, TotalTime:%v, TrackingId:%v\n",
					resp.StatusCode(), resp.Proto(), resp.Time(), resp.Header().Get(TrackingId))
			}

			ti := resp.Request.TraceInfo()
			fmt.Printf("Trace: %+v\n", ti)
		}()
	}
	wg.Wait()
}

func GetRestyClient() *resty.Client {
	// Create an HTTP/2 transport
	tr := &http2.Transport{}

	// Create an HTTP client with the transport
	httpClient := &http.Client{
		Transport: tr,
	}

	restyClient := resty.NewWithClient(httpClient)
	// Registering Request Middleware
	restyClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Now you have access to Client and current Request object
		// manipulate it as per your need

		// SetHeader(TrackingId, uuid.New().String()).
		req.SetHeader(TrackingId, uuid.New().String())

		return nil // if its success otherwise return error
	})

	// Registering Response Middleware
	restyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Now you have access to Client and current Response object
		// manipulate it as per your need

		return nil // if its success otherwise return error
	})

	return restyClient
}
