SHELL := /usr/bin/env bash

.PHONY: tidy envelope help

help:
	@printf '%s\n' \
		'make tidy                   Run go mod tidy for all modules' \
		'make release VERSION=<ver>  Tag and release all modules' \
		'make envelope                Regenerate envelope protobuf bindings'

tidy:
	@set -euo pipefail; \
	directory="$$(pwd)"; \
	modules=(); \
	while IFS= read -r module; do \
		modules+=("$$module"); \
	done < <(find "$$directory" -name go.mod -type f -not -path "*/.git/*" -print | sed 's|/go.mod$$||' | sort); \
	for module in "$${modules[@]}"; do \
		if [[ "$$module" == "$$directory" ]]; then \
			continue; \
		fi; \
		echo "tidy: $${module#$$directory/}"; \
		( \
			cd "$$module" && \
			go mod tidy \
		); \
	done; \
	if [[ -f "$$directory/go.mod" ]]; then \
		echo "tidy: ./"; \
		( \
			cd "$$directory" && \
			go mod tidy \
		); \
	fi


envelope:
	@set -euo pipefail; \
	mkdir -p envelope; \
	(cd idl/envelope && protoc --proto_path=. --go_out=paths=source_relative:../../envelope envelope.proto)
