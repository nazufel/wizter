# Dockerfile

# multi-stage Dockerfile to build the app

# first phase: build the binary with all of the build dependencies
FROM golang:1.16.5-alpine as builder

# install all build depdencies that don't come with the image
RUN apk add protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# copy over the source code to be built
WORKDIR /go/src/wizards
ADD ./cmd/client/main.go .
ADD ./wizard/ ./wizard/
ADD ./go.mod .
ADD ./go.sum .

# install dependent libraries
RUN go mod download

# build the protobuf definition
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto

# build the binary
ENV CGO_ENABLED=0 
ENV GOOS=linux
RUN go build -a -installsuffix cgo -o client .

# final phase: copy the built binary into a scratch image to run
FROM scratch

## set the working directory
WORKDIR /go/bin/wizards

## copy the combiled binary to the smaller image
COPY --from=builder /go/src/wizards/client .

## start the binary when the container starts
ENTRYPOINT ["./client"]