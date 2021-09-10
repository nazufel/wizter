all: images load 

run-client:
	$(clean_command)
	bazelisk run //cmd/client

client-push:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-client:v1 .

images: client server

load: load-client load-server

load-client:
	$(clean_command):
	kind load docker-image wizter/cmd/client:image

load-server:
	$(clean_command):
	kind load docker-image wizter/cmd/server:image

run-server:
	$(clean_command)
	bazelisk run //cmd/server

server-push:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-server:v1 .

proto: 
	$(clean_command)
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto