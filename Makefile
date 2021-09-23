all: images load-all

build-clients:
	$(clean_command)
	bazelisk run //cmd/client:v1
	bazelisk run //cmd/client:v2
	bazelisk run //cmd/client:v3

build-servers:
	$(clean_command)
	bazelisk run //cmd/server:v1
	bazelisk run //cmd/server:v2
	bazelisk run //cmd/server:v3

deploy-v1:
	$(clean_command)
	kubectl apply -f infra/wizards-namespace.yaml
	kubectl apply -f infra/mongo-pv.yaml
	kubectl apply -f infra/mongo-pvc.yaml
	kubectl apply -f infra/mongo-svc.yaml
	kubectl apply -f infra/mongo-deployment.yaml
	kubectl apply -f infra/wizards-client-configmap.yaml
	kubectl apply -f infra/wizards-server-configmap.yaml
	kubectl apply -f infra/wizards-server-svc.yaml
	kubectl apply -f infra/wizards-server-deployment-1.yaml
	kubectl apply -f infra/wizards-client-deployment-1.yaml

deploy-v2:
	$(clean_command)
	kubectl apply -f infra/wizards-namespace.yaml
	kubectl apply -f infra/mongo-pv.yaml
	kubectl apply -f infra/mongo-pvc.yaml
	kubectl apply -f infra/mongo-svc.yaml
	kubectl apply -f infra/mongo-deployment.yaml
	kubectl apply -f infra/wizards-client-configmap.yaml
	kubectl apply -f infra/wizards-server-configmap.yaml
	kubectl apply -f infra/wizards-server-svc.yaml
	kubectl apply -f infra/wizards-server-deployment-2.yaml
	kubectl apply -f infra/wizards-client-deployment-2.yaml

deploy-v3:
	$(clean_command)
	kubectl apply -f infra/wizards-namespace.yaml
	kubectl apply -f infra/mongo-pv.yaml
	kubectl apply -f infra/mongo-pvc.yaml
	kubectl apply -f infra/mongo-svc.yaml
	kubectl apply -f infra/mongo-deployment.yaml
	kubectl apply -f infra/wizards-client-configmap.yaml
	kubectl apply -f infra/wizards-server-configmap.yaml
	kubectl apply -f infra/wizards-server-svc.yaml
	kubectl apply -f infra/wizards-server-deployment-3.yaml
	kubectl apply -f infra/wizards-client-deployment-3.yaml

push-client:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-client:v1 .

images: build-clients build-servers

load-all: load-clients load-servers

load-clients:
	$(clean_command):
	kind load docker-image wizter/cmd/client:v1
	kind load docker-image wizter/cmd/client:v2
	kind load docker-image wizter/cmd/client:v3

load-servers:
	$(clean_command):
	kind load docker-image wizter/cmd/server:v1
	kind load docker-image wizter/cmd/server:v2
	kind load docker-image wizter/cmd/server:v3

run-client:
	$(clean_command)
	bazelisk run //cmd/client

run-server:
	$(clean_command)
	bazelisk run //cmd/server

proto: 
	$(clean_command)
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative wizard/wizard.proto

push-server:
	$(clean_command)
	docker build -f cmd/client/Dockerfile -t wizards-server:v1 .