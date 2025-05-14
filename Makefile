cli:
	docker compose up -d cdk

bash:
	docker compose run --rm cdk bash

test:
	docker compose run --rm cdk bash -c "go test ./lib/googlesheets/..."

deploy:
	docker compose run --rm cdk cdk deploy

clean:
	docker compose kill
	docker compose rm -f
