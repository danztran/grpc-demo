install:
	brew tap bufbuild/buf
	brew install buf
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go mod download

start-server:
	go run server/main.go

start-client:
	go run client/main.go

generate:
	buf generate
