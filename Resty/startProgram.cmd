set EVENTHUB_CONNECTION_STRING=
go run .\TestParallel.go .\TestResty.go .\HttpTrace.go -r 1 -url https://httpbin.org/get -eventHub %EVENTHUB_CONNECTION_STRING%