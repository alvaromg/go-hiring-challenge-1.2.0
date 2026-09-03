tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

run ::
	@go run cmd/server/main.go

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

test-integration ::
	@go test -v -count=1 -race ./test/... -coverprofile=integration_coverage.out -covermode=atomic -coverpkg github.com/mytheresa/go-hiring-challenge...
	@go tool cover -html=integration_coverage.out

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down

swag ::
	@swag init --dir ./cmd/server/,./infra/rest/,./infra/api --output ./swagger


