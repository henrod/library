build/protos:
	cd proto/library && buf push
	buf generate buf.build/henrod/library

run/api:
	go run app/service/main.go

lint:
	golangci-lint run
	# buf lint buf.build/henrod/library # this lint is incompatible with Google API design guide