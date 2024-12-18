default: all

all: mdbxv

.PHONY: mdbxv
mdbxv:
	@echo "Building MDBX Viewer..."
	@docker compose up -d --build