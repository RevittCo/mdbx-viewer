# DEPRECATED: This Makefile is deprecated in favour of Taskfile.yml (https://taskfile.dev).
# Run `task --list-all` to see available tasks. The Makefile remains as a thin
# wrapper for the transition period and will be removed in a future release.

define DEPRECATION_BANNER


╔══════════════════════════════════════════════════════════════════════╗
║                                                                      ║
║   WARNING: This Makefile is DEPRECATED.                              ║
║                                                                      ║
║   Use Taskfile (https://taskfile.dev) instead:                       ║
║       task --list-all      # show available tasks                    ║
║       task <task-name>     # run a task                              ║
║                                                                      ║
║   Install task:  brew install go-task                                ║
║   This Makefile will be removed in a future release.                 ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝


endef
export DEPRECATION_BANNER

$(info $(DEPRECATION_BANNER))

default: all

all: mdbxv

.PHONY: mdbxv
mdbxv:
	@echo "Building MDBX Viewer..."
	@docker compose up -d --build --force-recreate
