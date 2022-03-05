protoc:
	protoc --go_out=plugins=grpc:. proto/*/*.proto

server:
	go run cmd/server/main.go

server2:
	go run cmd/server_2/main.go

client:
	go run cmd/client/main.go

test:
	go test ./...
	