CORES := $$(getconf _NPROCESSORS_ONLN)

.PHONY: all
all: bootstrap package test

.PHONY: bootstrap
bootstrap:
	docker-compose build go

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

.PHONY: bbtest
bbtest:
	@echo "[info] stopping older runs"
	@(docker rm -f $$(docker-compose ps -q) 2> /dev/null || :) &> /dev/null
	@echo "[info] running bbtest"
	@docker-compose run --rm bbtest
	@echo "[info] stopping runs"
	@(docker rm -f $$(docker-compose ps -q) 2> /dev/null || :) &> /dev/null
	@(docker rm -f $$(docker ps -aqf "name=bbtest") || :) &> /dev/null

.PHONY: package
package:
	docker-compose run --rm package
	@#docker-compose build service

.PHONY: verify
verify:
	@echo "\nRAML 0.8"
	@find ./bin -type f -name "darwin*" -exec ./{} test spec/raml/v08/api.raml \;
	@echo "\nRAML 1.0"
	@find ./bin -type f -name "darwin*" -exec ./{} test spec/raml/v10/api.raml \;
	@echo "\nSWAGGER 2.0"
	@find ./bin -type f -name "darwin*" -exec ./{} test spec/swagger/v20/api.json \;
	@echo "\nSWAGGER 3.0"
	@find ./bin -type f -name "darwin*" -exec ./{} test spec/swagger/v30/api-with-examples.yaml \;

