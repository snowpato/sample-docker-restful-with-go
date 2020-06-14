module github.com/snowpato/sample-docker-restful-with-go/server

go 1.14

require google.golang.org/grpc v1.29.1

require (
	github.com/snowpato/sample-docker-restful-with-go/proto v0.0.0
	go.mongodb.org/mongo-driver v1.3.4
	google.golang.org/protobuf v1.24.0 // indirect
)

replace github.com/snowpato/sample-docker-restful-with-go/proto => ../proto
