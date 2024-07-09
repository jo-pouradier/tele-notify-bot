proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/metrics.proto && \
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/metrics.proto

server:
	go run ./cmd/master/main.go -tls

client:
	go run ./cmd/ping/main.go -tls
