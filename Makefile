build-controller:
	go build -o bin/controller controller/*.go

build-worker:
	go build -o bin/worker worker/*.go

build-proto:
	cd proto && \
	protoc --go_out=. --go_opt=module=thesis/proto --go-grpc_out=. --go-grpc_opt=module=thesis/proto worker.proto

clean:
	rm -rf bin/*
