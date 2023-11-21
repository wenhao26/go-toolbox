SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o ticker_example ticker_example.go
go build -ldflags "-s -w" -o ticker_example-ldflags ticker_example.go