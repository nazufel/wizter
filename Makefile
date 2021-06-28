all: images load 

client:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-client:v1 .

client-push:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-client:v1 .

images: client server

load: load-client load-server

load-client:
	$(clean_command):
	kind load docker-image wizards-client:v1

load-server:
	$(clean_command):
	kind load docker-image wizards-server:v1

server:
	$(clean_command)
	docker build -f cmd/server/Dockerfile -t wizards-server:v1 .

server-push:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-server:v1 .

proto: 
	$(clean_command)
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto