// build with go build -o tlsChecker.exe
// test: .\tlschecker.exe --url https://alanfeng-test-vnetenc.azure-api.net/internal-status-0123456789abcdef --maxTls TLS1.3
// go build -o tlsCheckerLib.so -buildmode=c-shared .
package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"flag"

	// 	"net/http"
	"github.com/go-resty/resty/v2"
)

func main() {

	// https://alanfeng-test-vnetenc.azure-api.net/internal-status-0123456789abcdef
	var targetUrl = flag.String("url", "https://google.com", "Provide URL to send requests to")
	var maxTlsVersion = flag.String("maxTls", "TLS1.3", "Max TLS version to use")
	var timeoutInSeconds = flag.Int("timeout", 10, "Timeout in seconds")
	flag.Parse()

	resp, err := SendRequest(*targetUrl, *maxTlsVersion, *timeoutInSeconds)
	fmt.Println("Calling GET:", *targetUrl)

	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode)
	fmt.Println("TLS version:", resp.tlsVersion)
	fmt.Println("Response Body:", resp.ResponseBody)
	fmt.Println("DNSLookupDuration:", resp.DNSLookupDuration)
	fmt.Println("ConnTime:", resp.ConnTime)
	fmt.Println("TCPConnTime:", resp.TCPConnTime)
	fmt.Println("TLSHandshake:", resp.TLSHandshake)
	fmt.Println("ServerTime:", resp.ServerTime)
}

//export SendRequest
func SendRequest(url string, maxTlsVersion string, timeoutInSeconds int) (*ResponseDetails, error) {
	var (
		conn *tls.Conn
		err  error
	)

	versions := map[uint16]string{
		tls.VersionTLS10: "TLS 1.0",
		tls.VersionTLS11: "TLS 1.1",
		tls.VersionTLS12: "TLS 1.2",
		tls.VersionTLS13: "TLS 1.3",
	}

	versionsFlag := map[string]uint16{
		"TLS1.0": tls.VersionTLS10,
		"TLS1.1": tls.VersionTLS11,
		"TLS1.2": tls.VersionTLS12,
		"TLS1.3": tls.VersionTLS13,
	}

	var tlsConfig *tls.Config = &tls.Config{
		MaxVersion: versionsFlag[maxTlsVersion],
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialTLS: func(network, addr string) (net.Conn, error) {
				conn, err = tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	restyClient := resty.NewWithClient(client)
	restyClient.SetTimeout(time.Duration(timeoutInSeconds) * time.Second)

	resp, err := restyClient.
		R().
		EnableTrace().
		Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	ti := resp.Request.TraceInfo()
	response := &ResponseDetails{
		Timestamp:         resp.Request.Time,
		url:               url,
		StatusCode:        resp.StatusCode(),
		ResponseBody:      resp.String(),
		tlsVersion:        versions[conn.ConnectionState().Version],
		DNSLookupDuration: ti.DNSLookup,
		ConnTime:          ti.ConnTime,
		TCPConnTime:       ti.TCPConnTime,
		TLSHandshake:      ti.TLSHandshake,
		ServerTime:        ti.ServerTime,
	}

	return response, nil
}

type ResponseDetails struct {
	Timestamp time.Time

	TrackingId string

	StatusCode int

	url string

	tlsVersion string

	ResponseBody string

	DNSLookupDuration time.Duration

	ConnTime time.Duration

	TCPConnTime time.Duration

	TLSHandshake time.Duration

	ServerTime time.Duration
}
