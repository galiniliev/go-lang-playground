package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

	"golang.org/x/net/http2"
)

const TrackingId = "x-ms-tracking-id"

var eventHubCtx context.Context

func TestParallel() {
	fmt.Println("TestParallel: start time:", time.Now().UTC())

	connStr := os.Getenv("EVENTHUB_CONNECTION_STRING")
	hub, err := eventhub.NewHubFromConnectionString(connStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	eventHubCtx = ctx
	defer cancel()

	var wg sync.WaitGroup
	// Create a Resty Client
	client := GetRestyClient()
	MakeRequest(client, url, hub)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			MakeRequest(client, url, hub)
		}()
	}
	wg.Wait()
}

func MakeRequest(client *resty.Client, url string, hub *eventhub.Hub) {
	var trackingId = uuid.New().String()
	var start = time.Now()
	resp, err := client.R().
		SetHeader(TrackingId, trackingId).
		Get(url)

	var duration = time.Now().Sub(start)
	//fmt.Printf("TrackingId: %v, Status Code: %v, Duration:%v\n", trackingId, resp.StatusCode(), duration)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		// defer resp.Body.Close()
		// fmt.Printf("Status Code: %v , Proto: %v, TotalTime:%v, TrackingId:%v\n",
		// 	resp.StatusCode(), resp.Proto(), resp.Time(), resp.Header().Get(TrackingId))
	}

	var httpTrace = HttpLogInfo{
		RequestUrl:    resp.Request.URL,
		TotalDuration: duration,
		StatusCode:    resp.StatusCode(),
		Protocol:      resp.Proto(),
		Timestamp:     resp.Request.Time,
		TrackingId:    trackingId,
		Trace:         resp.Request.TraceInfo(),
	}

	fmt.Printf("TrackingId: %v, Trace: %+v\n", trackingId, httpTrace.Trace)

	traceString, err := json.Marshal(httpTrace)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var logData = string(traceString)

	// send a single message into a random partition
	err = hub.Send(eventHubCtx, eventhub.NewEventFromString(logData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("%v, %v \n", time.Now().UTC(), logData)
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
