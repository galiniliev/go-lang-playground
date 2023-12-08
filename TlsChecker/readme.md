Console application to make an HTTP call to an endpoint and log TLS version, IP address and some more connection details..

```powershell
.\tlschecker.exe --url https://alanfeng-test-vnetenc.azure-api.net/ --maxTls TLS1.3 --timeout 10
Calling GET: https://alanfeng-test-vnetenc.azure-api.net/
Response Status Code: 404
TLS version: TLS 1.3
Remote Address: 20.47.146.143:443
Response Body: { "statusCode": 404, "message": "Resource not found" }
DNSLookupDuration: 0s
ConnTime: 286.368ms
TCPConnTime: 0s
TLSHandshake: 0s
ServerTime: 63.5124ms
```

to build locally 
```
go build -o tlsChecker.exe   
```

to test locally 
```
go run .\tlschecker.go --url https://alanfeng-test-vnetenc.azure-api.net/internal-status-0123456789abcdef --maxTls TLS1.3
```