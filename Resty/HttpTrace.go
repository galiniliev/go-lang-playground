package main

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type HttpLogInfo struct {
	RequestUrl string

	TotalDuration time.Duration

	StatusCode int

	TrackingId string

	Timestamp time.Time

	Protocol string

	Trace resty.TraceInfo
}
