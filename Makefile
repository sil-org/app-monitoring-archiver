cli:
	docker-compose up -d app

bash:
	docker-compose run --rm app bash

test:
	docker-compose run --rm app bash -c "go test ./lib/googlesheets/..."

clean:
	docker-compose kill
	docker-compose rm -f
