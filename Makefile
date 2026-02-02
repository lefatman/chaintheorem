# Minimal repo ergonomics. Real build/run targets land in later batches.

.PHONY: help fmt vet test proto

help:
	@echo "Targets:"
	@echo "  fmt    - gofmt ./..."
	@echo "  vet    - go vet ./..."
	@echo "  test   - go test ./..."
	@echo "  proto  - generate Go protobuf types (requires protoc + protoc-gen-go)"

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

proto:
	@./scripts/gen_proto.sh
