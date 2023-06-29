build-controller:
	go build -o bin/controller controller/*.go

build-worker:
	go build -o bin/worker worker/*.go

gen-proto:
	cd proto && \
	protoc --go_out=. --go_opt=module=thesis/proto --go-grpc_out=. --go-grpc_opt=module=thesis/proto worker-service.proto

clean:
	rm -rf bin/*

run-controller:
	go run controller/*.go

run-worker:
	go run worker/*.go
