
devenv:
	$(clean_command):
	export MONGO_HOST=mongo.default.svc.cluster.local
	export MONGO_USER=mongoadmin
	export MONGO_PASSWORD=admin123
	export MONGO_DATABASE=wizard

wizard: 
	$(clean_command)
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/wizard/wizard.proto