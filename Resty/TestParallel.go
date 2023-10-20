package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const TrackingId = "x-ms-tracking-id"
const RequestBatchSize = 300

var eventHubCtx context.Context
var eventHub *eventhub.Hub

func TestParallel(targetUrl string, numberOfRequests int, eventHubConnString string) {
	fmt.Println("TestParallel: start time:", time.Now().UTC())
	fmt.Printf("TestParallel: targetUrl:%v, numberOfRequests:%v eventHub:%v\n", targetUrl, numberOfRequests, eventHubConnString)

	if eventHubConnString != "" {
		hub, err := eventhub.NewHubFromConnectionString(eventHubConnString)
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		eventHubCtx = ctx
		eventHub = hub
		defer cancel()
	}

	// Create a Resty Client
	client := GetRestyClient()

	totalRequests := 0
	for {
		var wg sync.WaitGroup

		for i := 0; i < RequestBatchSize; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				MakeRequest(client, targetUrl)
			}()
		}
		wg.Wait()

		totalRequests += RequestBatchSize
		fmt.Printf("Time:%v Total requests executed:%v\n", time.Now().UTC(), totalRequests)
		trace := TaceEntry{
			Timestamp: time.Now().UTC(),
			Level:     4,
			Message:   fmt.Sprintf("Time:%v Total requests executed:%v\n", time.Now().UTC(), totalRequests),
			Properties: map[string]interface{}{
				"Batch":            RequestBatchSize,
				"TotalRequests":    totalRequests,
				"NumberOfRequests": numberOfRequests,
			},
		}

		traceString, err := json.Marshal(trace)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		var logData = string(traceString)
		LogEvent("Events", logData, "Events_mapping")

		time.Sleep(1 * time.Second)
		if numberOfRequests > 0 && totalRequests >= numberOfRequests {
			fmt.Printf("Exiting...")
			return
		}
	}
}

func MakeRequest(client *resty.Client, url string) {
	var trackingId = uuid.New().String()
	var start = time.Now()
	resp, err := client.R().
		SetHeader(TrackingId, trackingId).
		Get(url)

	var duration = time.Now().Sub(start)
	//fmt.Printf("TrackingId: %v, Status Code: %v, Duration:%v\n", trackingId, resp.StatusCode(), duration)

	if err != nil {
		LogError(trackingId, err)
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

	fmt.Printf("TrackingId: %v, Trace: %+v\n", trackingId, httpTrace)

	traceString, err := json.Marshal(httpTrace)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var logData = string(traceString)
	LogEvent("Requests", logData, "Requests_mapping")

	// fmt.Printf("%v, %v \n", time.Now().UTC(), logData)
}

func LogEvent(table string, logData string, mapping string) {
	// send a single message into a random partition
	if eventHub != nil {
		event := eventhub.NewEventFromString(logData)
		event.Properties = make(map[string]interface{})
		event.Properties["Table"] = table
		event.Properties["IngestionMappingReference"] = mapping

		err := eventHub.Send(eventHubCtx, event)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func LogError(trackingId string, err error) {
	trace := TaceEntry{
		Timestamp: time.Now().UTC(),
		Level:     2,
		Message:   fmt.Sprintf("Error:%v", err),
		Properties: map[string]interface{}{
			"err": err,
		},
	}

	traceString, err := json.Marshal(trace)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var logData = string(traceString)
	LogEvent("Events", logData, "Events_mapping")
}

func GetRestyClient() *resty.Client {
	restyClient := resty.New().
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
