SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o del_fcm_push_logs del_fcm_push_logs.go
go build -ldflags "-s -w" -o del_fcm_push_logs-ldflags del_fcm_push_logs.go