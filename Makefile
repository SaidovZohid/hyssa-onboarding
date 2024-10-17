CURRENT_DIR=$(shell pwd)

-include .env

PSQL_CONTAINER_NAME?=postgres-cli
PROJECT_NAME?=go-mono-repo
PSQL_URI?=postgres://postgres:1234@localhost:5432/${PROJECT_NAME}?sslmode=disable
	
TAG=latest


help: ## shows this help
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: sqlc
sqlc: ## sqlc generates code
	sqlc generate

.PHONY: createdb
createdb: ## creates a database
	docker exec -it ${PSQL_CONTAINER_NAME} createdb -U postgres ${PROJECT_NAME}

.PHONY: execdb
execdb: ## executes a database
	docker exec -it ${PSQL_CONTAINER_NAME} psql -U postgres ${PROJECT_NAME}

.PHONY: dropdb
dropdb: ## drops a database
	docker exec -it ${PSQL_CONTAINER_NAME} dropdb -U postgres ${PROJECT_NAME}

.PHONY: execdb
cleandb: ## cleans a database
	docker exec -it ${PSQL_CONTAINER_NAME} psql -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" ${PSQL_URI}

.PHONY: migrate_up
migrate_up: ## migrates up
	goose -dir migrate/migrations postgres "${PSQL_URI}" up

.PHONY: migrate_down
migrate_down: ## migrates down
	goose -dir migrate/migrations postgres "${PSQL_URI}" down

.PHONY: migrate_status
migrate_status: ## migrates status
	goose -dir migrate/migrations postgres "${PSQL_URI}" status

.PHONY: migrate_create
migrate_create: ## migrates create
	goose -s -dir migrate/migrations create ${NAME} sql

build_image: ## builds the image
	docker build --rm -t "${REGISTRY_NAME}/${PROJECT_NAME}:${TAG}" .

push_image: ## pushes the image
	docker push "${REGISTRY_NAME}/${PROJECT_NAME}:${TAG}"

proto: ## generates the proto files
	rm -f generated/**/*.go
	rm -f doc/swagger/*.swagger.json
	mkdir -p generated
	protoc \
		--proto_path=protos --go_out=generated --go_opt=paths=source_relative \
		--go-grpc_out=generated --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=generated --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=swagger_docs,use_allof_for_refs=true,disable_service_tags=true \
			protos/**/*.proto
