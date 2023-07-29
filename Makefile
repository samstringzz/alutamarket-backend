include .env
export

postgresinit:
	docker run --name postgresaluta -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:15-alpine

postgres:
	docker exec -it postgresaluta psql

createdb:
	docker exec -it postgresaluta createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it postgresaluta dropdb $(DB_NAME)

migration: 
	migrate create -ext sql -dir db/migrations add_transaction_table

migrationup:
	# migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose force 20230704132640 up 
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up 


migrationdown:
	# migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose force 20230705085541 down 
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose  down 

.PHONY: postgresinit postgres createdb dropdb migration migrationdown migrationup
