# Makefile for proto file generation
PROTO_DIR = dsl
TPL_DIR = $(PROTO_DIR)/tpl
OUTPUT = $(PROTO_DIR)/demo_service.proto

.PHONY: clean build
code:
	goctl rpc protoc dsl/demo_service.proto --go_out=./  --go-grpc_out=./  --zrpc_out=./ --client=true --multiple=true --style=go_zero




up:
	#go get gitlab.leadsea.cn/go-template/base-common
	#go mod tidy
run:
	go run .\rpc_demo.go



build: clean
	@echo "Generating merged proto file..."
	@mkdir -p $(PROTO_DIR)
	@echo 'syntax = "proto3";' > $(OUTPUT)
	@echo 'package rpc;' >> $(OUTPUT)
	@echo 'option go_package ="./rpc";' >> $(OUTPUT)
	@for file in $$(find $(TPL_DIR) -type f -name '*.proto'); do \
		echo "Processing $$file"; \
		sed '/syntax = "proto3";/d;/import "tpl\/base.proto";/d;/package rpc;/d;/option go_package = ".\/rpc";/d' $$file >> $(OUTPUT); \
	done

clean:
	@echo "Cleaning generated files..."
	@rm -f $(OUTPUT)