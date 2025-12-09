cli:
	docker compose up -d cdk

bash:
	docker compose run --rm cdk bash

test:
	go test ./...

deploy:
	docker compose run --rm cdk cdk deploy

clean:
	docker compose kill
	docker compose rm -f

update:
	go get -u ./... && go mod tidy
	cd cdk && go get -u ./... && go mod tidy && cd ..
