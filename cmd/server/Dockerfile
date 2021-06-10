# Dockerfile

# Dockerfile to build the server

# first phase: build with all of the build dependencies
FROM golang:1.16.5-alpine as builder

RUN apk add protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# set the
WORKDIR /go/src/wizards
ADD ./cmd/server/main.go .
ADD ./wizard/ ./wizard/
ADD ./go.mod .
ADD ./go.sum .

RUN go mod download
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto

ENV CGO_ENABLED=0 
ENV GOOS=linux
RUN go build -a -installsuffix cgo -o server .

# final phase: copy binary to scratch container
FROM alpine

## set the working directory
WORKDIR /go/bin/wizards

## copy the combiled binary to the smaller image
COPY --from=builder /go/src/wizards/server .

## start the binary when the container starts
ENTRYPOINT ["./server"]

# CMD [ "/bin/ash", "-c", "sleep 111111111" ]
