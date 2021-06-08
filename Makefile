client:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t nazufel/telepresence-demo-client:v1 .

devenv:
	$(clean_command):
	export MONGO_HOST=mongo.default.svc.cluster.local
	export MONGO_USER=mongoadmin
	export MONGO_PASSWORD=admin123
	export MONGO_DATABASE=wizard

server:
	$(clean_command)
	docker build -f cmd/server/Dockerfile -t nazufel/telepresence-demo-server:v1 .a

wizard: 
	$(clean_command)
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto