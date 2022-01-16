build/protos:
	buf generate buf.build/henrod/library

run/api:
	go run app/service/main.go