package main

// Import resty into your code and refer it as `resty`.
import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const url = "https://apim-mckq3zpiubjie.azure-api.net/mock/get"
const url1 = "https://httpbin.org/get"

var eventHubConnectionString string

func main() {

	var requests = flag.Int("requests", 10, "Provide number of requests to send")
	var targetUrl = flag.String("url", url, "Provide to send requests to")
	var eventHubConnStrPtr = flag.String("eventHub", "", "Provide connection string for event hub")
	eventHubConnectionString = *eventHubConnStrPtr

	requestsEnv, err := strconv.Atoi(os.Getenv("load-test-requests"))
	if err == nil {
		requests = &requestsEnv
	}

	var urlEnv = os.Getenv("load-test-url")
	if urlEnv != "" {
		targetUrl = &urlEnv
	}

	var eventHubEnv = os.Getenv("load-test-eventHub")
	if urlEnv != "" {
		eventHubConnectionString = eventHubEnv
	}
	flag.Parse()

	fmt.Printf("Received flags requests:%v targetUrl:%v\n", *requests, *targetUrl)

	// TestSingleGet()
	TestParallel(*targetUrl, *requests)
}

func TestSingleGet() {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(url)

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Printf("  Trace Info    :%+v\n", ti)
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())

	/* Output
	   Response Info:
	     Error      : <nil>
	     Status Code: 200
	     Status     : 200 OK
	     Proto      : HTTP/2.0
	     Time       : 457.034718ms
	     Received At: 2020-09-14 15:35:29.784681 -0700 PDT m=+0.458137045
	     Body       :
	     {
	       "args": {},
	       "headers": {
	         "Accept-Encoding": "gzip",
	         "Host": "httpbin.org",
	         "User-Agent": "go-resty/2.4.0 (https://github.com/go-resty/resty)",
	         "X-Amzn-Trace-Id": "Root=1-5f5ff031-000ff6292204aa6898e4de49"
	       },
	       "origin": "0.0.0.0",
	       "url": "https://httpbin.org/get"
	     }

	   Request Trace Info:
	     DNSLookup     : 4.074657ms
	     ConnTime      : 381.709936ms
	     TCPConnTime   : 77.428048ms
	     TLSHandshake  : 299.623597ms
	     ServerTime    : 75.414703ms
	     ResponseTime  : 79.337µs
	     TotalTime     : 457.034718ms
	     IsConnReused  : false
	     IsConnWasIdle : false
	     ConnIdleTime  : 0s
	     RequestAttempt: 1
	     RemoteAddr    : 3.221.81.55:443
	*/

}
