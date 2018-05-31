CORES := $$(getconf _NPROCESSORS_ONLN)

.PHONY: all
all: bootstrap package test

.PHONY: bootstrap
bootstrap:
	docker-compose build go
	docker-compose run --rm sync

.PHONY: fetch
fetch:
	docker-compose run fetch

.PHONY: build-lint
build-lint:
	docker-compose build lint

.PHONY: build-sync
build-sync:
	docker-compose build sync

.PHONY: build-package
build-package:
	docker-compose build package

.PHONY: lint
lint:
	docker-compose run --rm lint || :

.PHONY: sync
sync:
	docker-compose run --rm sync

.PHONY: test
test:
	docker-compose run --rm test

.PHONY: bench
bench:
	docker-compose run --rm bench

.PHONY: package
package:
	docker-compose run --rm package
	@#docker-compose build service

.PHONY: verify
verify:
	@eval $(eval ct=$(shell sh -c 'find ./bin -type f -name "darwin*"' | awk '{print $$1}'))

	@echo "[info] test all"
	@find ./spec -type f -name "*api*" -exec $(ct) test {} \;

.PHONY: bbtest
bbtest:
	@eval $(eval ct=$(shell sh -c 'find ./bin -type f -name "darwin*"' | awk '{print $$1}'))

	@echo "[info] test bbtest"
	@$(ct) test bbtest/raml/api.raml
