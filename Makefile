.PHONY: rebuild
rebuild:
	docker-compose -f docker-compose.yml down
	docker volume rm p2p_data
	docker volume create --name=p2p_data
	docker-compose -f docker-compose.yml build --force-rm
	docker-compose -f docker-compose.yml up -d
	docker image prune -f

.PHONY: build
build:
	@echo "-- building binary"
	go build \
		-o ./bin/p2pcrawler \
		./cmd/p2pcrawler
