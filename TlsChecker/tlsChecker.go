package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"flag"

	// 	"net/http"
	"github.com/go-resty/resty/v2"
)

func main() {
	var (
		conn *tls.Conn
		err  error
	)

	// https://alanfeng-test-vnetenc.azure-api.net/internal-status-0123456789abcdef
	var targetUrl = flag.String("url", "https://google.com", "Provide URL to send requests to")
	flag.Parse()

	tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig

	client := &http.Client{
		Transport: &http.Transport{
			DialTLS: func(network, addr string) (net.Conn, error) {
				conn, err = tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	restyClient := resty.NewWithClient(client)

	resp, err := restyClient.R().Get(*targetUrl)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	versions := map[uint16]string{
		tls.VersionTLS10: "TLS 1.0",
		tls.VersionTLS11: "TLS 1.1",
		tls.VersionTLS12: "TLS 1.2",
		tls.VersionTLS13: "TLS 1.3",
	}

	fmt.Println("Calling GET:", *targetUrl)
	fmt.Println("Response Status Code:", resp.StatusCode())
	fmt.Println("TLS version:", versions[conn.ConnectionState().Version])
	fmt.Println("Response Body:", resp)
}
