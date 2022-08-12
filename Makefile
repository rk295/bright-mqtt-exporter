.PHONY: run
run:
	./run

build:
	cd cmd/bright-mqtt-exporter && go build -o ../../bright-mqtt-exporter

.PHONY: docker
docker:
	docker-compose build

.PHONY: restart
restart:
	docker-compose up -d

.PHONY: logs
logs:
	docker-compose logs -f --tail 100
