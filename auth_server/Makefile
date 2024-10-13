LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml


install-deps:
	@echo ">  Install dependencies..."
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	@echo ">  Get dependencies..."
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	

generate:
	@echo ">  Generate User API..."
	make generate-user-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto

docker-build-and-push:
# Сборка образа из докерфайла под целевую платформу --platform с тегом -t <PATH_TO_REGESTRY/app_name:version <PATH_TO_SOURCES>
# Образ будет положен локально
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/algol/auth-server:v0.0.1 .
# Подключение к удаленному regestry по логину -u и паролю -p
	docker login -u token -p CRgAAAAATQUDNnaAfiJADdGp2n18L8EDg8xkHJbH cr.selcloud.ru/algol
# Загрузка образа докера в удаленный реестр
	docker push cr.selcloud.ru/algol/auth-server:v0.0.1 

docker-migration-push:
	@echo ">  dgdfgdfhfh..."
	docker tag postgres:14-alpine3.17 cr.selcloud.ru/algol/postgres:14-alpine3.17
	docker login -u token -p CRgAAAAATQUDNnaAfiJADdGp2n18L8EDg8xkHJbH cr.selcloud.ru/algol
	docker push cr.selcloud.ru/algol/postgres:14-alpine3.17
