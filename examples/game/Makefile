PROTO_DIR := proto
PB_DIR := internal/pb

.PHONY: proto
proto:
	@echo "==> Proto: generating Go code"
	@if [ -z "$$(ls -1 $(PROTO_DIR)/*.proto 2>/dev/null)" ]; then \
		echo "Error: no .proto files found in $(PROTO_DIR)"; \
		exit 1; \
	fi
	@echo "Where: $(PROTO_DIR) --> $(PB_DIR)"
	@files=$$(ls -1 $(PROTO_DIR)/*.proto | xargs -n1 basename | paste -sd ", " -); \
		echo "Files: $$files"
	@if ! docker run --rm \
		-v "$(CURDIR)/$(PROTO_DIR):/defs" \
		-v "$(CURDIR)/$(PB_DIR):/out" \
		rvolosatovs/protoc \
		--proto_path=/defs \
		--go_out=/out \
		--go_opt=paths=source_relative \
		/defs/*.proto; then \
		echo "Error: generation failed. Check Docker is running and you can access the Docker socket."; \
		exit 1; \
	fi
	@echo "---------------------------------"
	@count=$$(ls -1 $(PB_DIR)/*.pb.go 2>/dev/null | wc -l | tr -d ' '); \
		echo "OK: generated $$count file(s)"
	@ls -1 $(PB_DIR)/*.pb.go 2>/dev/null | sed 's|^| - |'
