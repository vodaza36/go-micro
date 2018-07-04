build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/vodaza36/go-micro/user-service \
		proto/user/user.proto
	docker build -t user-service .

run:
	docker run -p 50051:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns user-service