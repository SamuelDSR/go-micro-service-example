module github.com/samueldsr/shippy-cli-consignment

go 1.14

require (
	github.com/samueldsr/shippy-service-consignment v0.0.0
	google.golang.org/grpc v1.30.0
)

replace github.com/samueldsr/shippy-service-consignment => /home/samuel_vita/Documents/go-micro-service/shippy/shippy-service-consignment
