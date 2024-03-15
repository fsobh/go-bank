IMAGE_NAME = 'simple-bank'
USER = root
PASSWORD = secret
PORT = 5432
DB_NAME = simple_bank

createpostgres:
	docker run --name $(IMAGE_NAME) -e POSTGRES_USER=$(USER) -e POSTGRES_PASSWORD=$(PASSWORD) -p $(PORT):$(PORT) -d postgres

createdb:
	docker exec -it $(IMAGE_NAME) createdb --username=$(USER) --owner=$(USER) $(DB_NAME)

dropdb:
	docker exec -it $(IMAGE_NAME) dropdb $(DB_NAME)

stoppostgres:
	docker stop $(IMAGE_NAME)

runpostgres:
	docker start $(IMAGE_NAME)

deletepostgres:
	docker rm $(IMAGE_NAME)

createmigrate:
	migrate create -ext sql -dir db/migration -seq init_schema # only run if the db/migration folder does not exist

migrateup:
	migrate -path db/migration -database "postgresql://$(USER):$(PASSWORD)@localhost:$(PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://$(USER):$(PASSWORD)@localhost:$(PORT)/$(DB_NAME)?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://$(USER):$(PASSWORD)@localhost:$(PORT)/$(DB_NAME)?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://$(USER):$(PASSWORD)@localhost:$(PORT)/$(DB_NAME)?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v $(PWD):/src -w /src kjconroy/sqlc generate

test:
	go test -v -coverprofile=testCoverage.out ./...
	go tool cover -html=testCoverage.out
server:
	go run main.go

mockdb:
	mockgen --package mockdb --destination db/mock/store.go github.com/fsobh/simplebank/db/sqlc Store



.PHONY: createpostgres createdb dropdb stoppostgres runpostgres deletepostgres createmigrate migrateup migratedown sqlc test server mockdb migratedown1 migrateup1
