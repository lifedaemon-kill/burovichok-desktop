
.PHONY: all
all: help

.PHONY: help
help:
	@echo "Available commands:"
	@echo "gen    -- generating proto files in go"

.PHONY: gen
gen:
	@echo "run proto generation..."
	@protoc -I proto pkg/proto/api.proto --go_out=internal/pb/gen/ --go-grpc_out=internal/pb/gen/
	@echo "end proto generation"