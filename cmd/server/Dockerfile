# Dockerfile

# Dockerfile to build the server


# first phase: build with all of the build dependencies
FROM golang:1.16.5-alpine as builder

RUN apk add protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# set the
WORKDIR /go/src
ADD ./cmd/server/main.go .
ADD ./wizard/ ./wizard/
ADD ./go.mod .
ADD ./go.sum .

RUN go mod download
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto

RUN go build -o server main.go

# final phase: copy binary to scratch container
# FROM ubuntu
# ENV MONGO_HOST "mongo.default.svc.cluster.local" 
# ENV MONGO_DATABASE "wizard" 
# ENV MONGO_PORT "27017" 
# ENV GRPC_PORT "9999"

## set the working directory
# WORKDIR /go/bin/server

## copy the combiled binary to the smaller image
# COPY --from=builder /go/src/server .

# RUN cp /go/src/server /go/bin/server

## start the binary when the container starts
# ENTRYPOINT ["./server"]

CMD [ "/bin/ash", "-c", "sleep 111111111" ]
