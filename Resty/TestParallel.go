package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"golang.org/x/net/http2"
)

const TrackingId = "x-ms-tracking-id"

func TestParallel() {

	fmt.Println("TestParallel: start time:", time.Now().UTC())

	var wg sync.WaitGroup
	// Create a Resty Client
	client := GetRestyClient()
	MakeRequest(client, url)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			MakeRequest(client, url)
		}()
	}
	wg.Wait()
}

func MakeRequest(client *resty.Client, url string) {
	var trackingId = uuid.New().String()
	var start = time.Now()
	resp, err := client.R().
		SetHeader(TrackingId, trackingId).
		Get(url)

	var duration = time.Now().Sub(start)
	fmt.Printf("TrackingId: %v, Status Code: %v, Duration:%v\n", trackingId, resp.StatusCode(), duration)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// defer resp.Body.Close()
		// fmt.Printf("Status Code: %v , Proto: %v, TotalTime:%v, TrackingId:%v\n",
		// 	resp.StatusCode(), resp.Proto(), resp.Time(), resp.Header().Get(TrackingId))
	}

	ti := resp.Request.TraceInfo()
	fmt.Printf("TrackingId: %v, Trace: %+v\n", trackingId, ti)
}

func GetRestyClient() *resty.Client {
	// Create an HTTP/2 transport
	tr := &http2.Transport{}

	// Create an HTTP client with the transport
	httpClient := &http.Client{
		Transport: tr,
	}

	restyClient := resty.NewWithClient(httpClient).
		EnableTrace()

	// Registering Request Middleware
	restyClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Now you have access to Client and current Request object
		// manipulate it as per your need

		// var trackingId = uuid.New().String()
		// req.SetHeader(TrackingId, trackingId)
		var trackingId = req.Header.Get(TrackingId)
		req.SetHeader("User-Agent", fmt.Sprintf("Resty/Tracking:%v", trackingId))

		return nil // if its success otherwise return error
	})

	// Registering Response Middleware
	restyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Now you have access to Client and current Response object
		// manipulate it as per your need

		// var trackingId = resp.Header().Get(TrackingId)
		// var duration = time.Now().Sub(resp.Request.Time)
		// fmt.Printf("OnAfterResponse - TrackingId: %v, Duration:%v\n", trackingId, duration)

		return nil // if its success otherwise return error
	})

	return restyClient
}
