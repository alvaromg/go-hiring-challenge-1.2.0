tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

run ::
	@go run cmd/server/main.go

test ::
	@go test -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

test-integration ::
	@go test -count=1 -race ./test/... -coverprofile=integration.out -covermode=atomic -coverpkg github.com/mytheresa/go-hiring-challenge...
	@go tool cover -html=integration.out


docker-up ::
	docker compose up -d

docker-down ::
	docker compose down


