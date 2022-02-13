build/protos:
	cd proto/library && buf push
	buf generate buf.build/henrod/library

run/api:
	go run app/service/main.go

lint:
	golangci-lint run
	# buf lint buf.build/henrod/library # this lint is incompatible with Google API design guide

deps:
	docker run --name library -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_PASSWORD=password -d --rm postgres