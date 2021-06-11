module github.com/nazufel/telepresence-demo

go 1.16

replace github.com/nazufel/telepresence-demo/wizard => ./wizard

require (
	github.com/enescakir/emoji v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.5.3
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)
