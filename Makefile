CORES := $$(getconf _NPROCESSORS_ONLN)

.PHONY: all
all: bootstrap package test bbtest

.PHONY: bootstrap
bootstrap:
	docker-compose build go
	docker-compose run --rm sync
	docker-compose run --rm package -t darwin,linux

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

	@echo "[info] test local"
	@find ./spec -type f -name "*api*" -exec $(ct) test {} \;

	@echo "[info] test remote"
	@$(ct) test http://petstore.swagger.io/v2/swagger.json

.PHONY: bbtest
bbtest:
	@echo "[info] stopping older runs"
	@(docker rm -f $$(docker-compose ps -q) 2> /dev/null || :) &> /dev/null
	@echo "[info] running bbtest"
	@docker-compose run --rm bbtest
	@echo "[info] stopping runs"
	@(docker rm -f $$(docker-compose ps -q) 2> /dev/null || :) &> /dev/null
	@(docker rm -f $$(docker ps -aqf "name=bbtest") || :) &> /dev/null

.PHONY: tracer
tracer:
	@eval $(eval ct=$(shell sh -c 'find ./bin -type f -name "darwin*"' | awk '{print $$1}'))
	#docker-compose run --rm --service-ports testee
	@echo "[info] test bbtest"
	@ \
		VERSION=1 \
		PORT=8080 \
		X="xValue" \
		SUPER_VALUE="filledSuperValue" \
		$(ct) test bbtest/raml/api.raml
