.PHONY: all build run clean

all: build run

build:
	docker-compose build

run:
	docker-compose -f docker-compose.yml up

run-d:
	docker-compose -f docker-compose.yml up -d

clean:
	docker-compose -f docker-compose.yml down
