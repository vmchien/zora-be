#GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION="0.1.1"
NAME="user-api"
PROJECT_DOCS_DIR=internal/server/project-docs
DOCS_SOURCE_DIR=../../docs


KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)
ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find ./internal -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find ../../api -name *.proto")
	V1_PROTO_FILES=$(shell $(Git_Bash) -c "find ../../api -name *.proto")
else
	INTERNAL_PROTO_FILES=$(shell find ./internal -name *.proto)
	API_PROTO_FILES=$(shell find ../../api -name *.proto)
	V1_PROTO_FILES=$(shell find ../../api -name *.proto)
endif


.PHONY: echo
echo:
	@echo "INTERNAL_PROTO_FILES=$(INTERNAL_PROTO_FILES)"
	@echo "API_PROTO_FILES=$(API_PROTO_FILES)"
	@echo "V1_PROTO_FILES=$(V1_PROTO_FILES)"

.PHONY: init
# init env
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest


.PHONY: config
# generate internal proto
config:
	@echo "Cleaning old generated files..."
	@find ./internal -type f \( -name "*.pb.go" \) -delete
	@echo "Generating files..."
	protoc --proto_path=./internal \
	       --proto_path=../../third_party \
 	       --go_out=paths=source_relative:./internal \
 	       --validate_out=paths=source_relative,lang=go:./internal \
	        $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	@echo "Cleaning old generated files..."
	@find ../../api -type f \( -name "*.pb.go" -o -name "*.pb.validate.go" -o -name "*.swagger.json" -o -name "*.openapi.json" -o -name "*_gen.proto" \) -delete
	@echo "Generating files..."
	protoc --proto_path=../../api \
               --proto_path=../../third_party \
               --go_out=paths=source_relative:../../api \
               --go-http_out=paths=source_relative:../../api \
               --go-grpc_out=paths=source_relative:../../api \
               --openapi_out=fq_schema_naming=true,default_response=false:../../api \
               --validate_out=paths=source_relative,lang=go:../../api \
               $(API_PROTO_FILES)

.PHONY: schema
# generate ent_schema for smp module
schema:
	@echo "ent schema generating......."

	@mkdir -p ./internal/data/ent
	@mkdir -p ./internal/data/ent_hook
	@mkdir -p ./internal/data/ent_interceptor
	@mkdir -p ./internal/data/gen_sql

	@for mod in user \
				; do \
		echo ">>>>>>> running for module: $$mod"; \
		echo "cleaning old files...."; \
		find "./internal/data/ent" -mindepth 1 -not -name 'init.go' -exec rm -rf {} +; \
		find "./internal/data/ent_hook" -type f \( -name '*_gen.go' \) -delete; \
		find "./internal/data/ent_interceptor" -type f \( -name '*_gen.go' \) -delete; \
		find "./internal/data/gen_sql" -type f \( -name '*_gen.sql' \) -delete; \
	done
	@go run ../../pkg/data/generate/generate.go user;

.PHONY: project-docs
# copy project documents for embedded dev docs UI
project-docs:
	@if [ -d "$(DOCS_SOURCE_DIR)" ]; then \
		mkdir -p $(PROJECT_DOCS_DIR); \
		rm -f $(PROJECT_DOCS_DIR)/apps-services-detail.html $(PROJECT_DOCS_DIR)/system-flow-redesign.html; \
		rm -rf $(PROJECT_DOCS_DIR)/assets; \
		cp $(DOCS_SOURCE_DIR)/apps-services-detail.html $(PROJECT_DOCS_DIR)/; \
		cp $(DOCS_SOURCE_DIR)/system-flow-redesign.html $(PROJECT_DOCS_DIR)/; \
		cp -R $(DOCS_SOURCE_DIR)/assets $(PROJECT_DOCS_DIR)/assets; \
	else \
		echo "Docs source directory not found, skipping project-docs"; \
	fi

.PHONY: swagger2
# generate swagger v2 file from proto
swagger2: project-docs
	mkdir -p internal/server/swagger-ui
	protoc --proto_path=../../api \
		--proto_path=../../third_party \
		--openapiv2_out internal/server/swagger-ui \
		--openapiv2_opt logtostderr=true,use_go_templates=true,json_names_for_fields=true,allow_merge=true,merge_file_name=gateway,allow_delete_body=true \
		$(V1_PROTO_FILES)
	cp internal/server/swagger-ui/gateway.swagger.json internal/server/swagger-ui/swagger.json
	rm -f internal/server/swagger-ui/gateway.swagger.json
	cp ../../templates/swagger-ui/*  internal/server/swagger-ui/


.PHONY: swagger3
# generate swagger v3 file from proto
swagger3: project-docs
	mkdir -p internal/server/swagger-ui/docs
	protoc --proto_path=../../api \
		--proto_path=../../third_party \
		--openapi_out=fq_schema_naming=true,default_response=false:internal/server/swagger-ui/docs \
		$(V1_PROTO_FILES)
	cp internal/server/swagger-ui/docs/openapi.yaml internal/server/swagger-ui/swagger.json
	rm -f internal/server/swagger-ui/docs/openapi.yaml
	cp ../../templates/swagger-ui/*  internal/server/swagger-ui/


.PHONY: build
# build
build: project-docs
	mkdir -p bin && cd cmd && GOEXPERIMENT=jsonv2,goroutineleakprofile go build -o ../bin/user -ldflags "-X main.Version=$(VERSION) -X main.Name=$(NAME)"

.PHONY: buildl
# build linux
buildl: project-docs
	mkdir -p bin && cd cmd && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOEXPERIMENT=jsonv2,goroutineleakprofile go build -o ../bin/user -trimpath -buildvcs=false -ldflags="-s -w -X main.Version=$(VERSION) -X main.Name=$(NAME)"

.PHONY: run
# run
run: project-docs
	cd cmd/ && kratos run

.PHONY: wire
# DI with google wire
WIRE_VERSION := v0.7.0
wire:
	cd cmd/ && go mod tidy && go run github.com/google/wire/cmd/wire@$(WIRE_VERSION)

.PHONY: error-codes
# generate error codes from yaml files for type-safe error handling
error-codes:
	cd ../../pkg/errors/gen && go run generate.go

.PHONY: all
# generate all
all:
	@echo "Name=${Name}"
	make config;
	make api;
	#make error-codes;
	make schema;
	#make swagger2;
	#make swagger3;
	make wire;
	make build;

.PHONY: docker
# generate all
docker:
	@echo "Name=${Name}, Version=${Version}"
	make config;
	make api;
	make error-codes;
	#make swagger2;
	make swagger3;
	make wire;
	make buildl;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
