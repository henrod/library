build/protos:
	cd proto/library && buf push
	buf generate buf.build/henrod/library

run/api:
	go run app/service/main.go

lint:
	golangci-lint run