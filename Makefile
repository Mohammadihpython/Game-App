
present-protobuf-golang:
		protoc --go_out=. --go-grpc_out=.   contract/protobuf/presence/presence.proto
start-docker-services:
	 docker compose up -d
 
Presence-service:
		go run cmd/presenceserver/server.go
Game-server:
	go run main.go