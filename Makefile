.PHONY: all clean ui generate build install test up down logs

all: build

clean:
	@rm -f wallawire
	@rm -rf ui/public
	@rm -f public

ui:
	@cd ui && yarn install && yarn build

generate: ui
	@ln -f -s ui/public
	@go generate

build: generate
	@./scripts/build.sh

install: build
	@mv wallawire $(GOBIN)

test:
	go test ./...

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f
