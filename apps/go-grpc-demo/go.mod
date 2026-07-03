module go-grpc-demo

go 1.25.1

require (
	dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/grpc v0.0.0
	dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger v0.0.0
)

require (
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/grpc v1.79.3 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/grpc => ../../go-packages/grpc

replace dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger => ../../go-packages/logger
