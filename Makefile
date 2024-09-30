install:
	brew install buf
	go install buf github.com/bufbuild/buf/cmd/buf@v1.14.0
	go install protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
	go install protoc-gen-openapiv2 github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1
	go install protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	go install protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install protoc-gen-validate github.com/envoyproxy/protoc-gen-validate@v1.0.2

start-server:
	go run server/main.go

start-client:
	go run client/main.go

generate:
	buf generate --template ./pb/buf.gen.yaml
