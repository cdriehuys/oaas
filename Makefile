.PHONY: container
container:
	docker build -t oaas:latest .

HOST_PORT ?= 8000

.PHONY: run
run: container
	docker run --publish $(HOST_PORT):8000 oaas:latest

BUILD_TAGS ?=

.PHONY: build-backend
build-backend:
	CGO_ENABLED=0 go build -gcflags="all=-N -l" -o ./build/main $(if $(BUILD_TAGS),-tags=$(BUILD_TAGS),) .

.PHONY: dev
dev:
	npx mprocs "make dev-ui" "make dev-backend"

.PHONY: dev-backend
dev-backend:
	go tool air

.PHONY: dev-ui
dev-ui:
	cd frontend; npm run dev -- --host

# By default the application starts without waiting for a debugger. To wait for
# a debugger to be attached before starting, set `WAIT_FOR_DEBUGGER` to a
# non-empty value in your `make` invocation.
WAIT_FOR_DEBUGGER ?=
DELVE_OPTS := $(if $(WAIT_FOR_DEBUGGER),,--continue)

.PHONY: run-backend
run-backend:
	# Subsequent reloads cause delve to still be running and bound to port 2345.
	# A short delay lets it exit between runs.
	sleep 1
	dlv exec ./build/main --listen=127.0.0.1:2345 --headless=true --api-version=2 --accept-multiclient $(DELVE_OPTS)

.PHONY: api-server
api-server:
	go generate ./api/internal/server
