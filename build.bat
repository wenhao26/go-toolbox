SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o filebeat_test logs_demo/main.go
;go build -ldflags "-s -w" -o ticker_example-ldflags ticker_example.go