VERSION_WIRE = v0.7.0
#VERSION_PROTOC :=


.PHONY: init
# init all necessary tools
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    && go install github.com/google/wire/cmd/wire@v0.7.0 \
    && go install github.com/envoyproxy/protoc-gen-validate@latest \
    && go install github.com/go-kratos/kratos/cmd/kratos/v2@latest \
    && go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest \
    && go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest \
    && go install entgo.io/ent/cmd/ent@latest

.PHONY: clean
# clean all installed binaries and cached files
clean:
	@echo "Cleaning installed binaries and cached files..."
	go clean -cache -modcache -i -r


.PHONY: update
# update all module dependencies
update:
	@echo "Updating dependencies in all modules..."
	@for mod in pkg \
				api/gateway \
				api/zalo \
				apps/gateway \
				apps/zalo \
        	; do \
    	echo "Running go get -u in $$mod"; \
    	(cd $$mod && go mod tidy && go get -u ./...); \
    done
	@echo "Dependencies updated."

.PHONY: build
# build all modules
build:
	@echo "Building all modules..."
	@for mod in apps/gateway \
				apps/zalo \
        	; do \
    	echo "Build module: $$mod"; \
    	(cd $$mod && make all); \
    done
	@echo "All modules built."
