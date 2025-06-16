start:
	sudo docker start postgres-container
postgres: 
	sudo docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres
createdb:
	sudo docker exec -it postgres-container createdb --username=root --owner=root ocr-database

dropdb:
	sudo docker exec -it postgres-container dropdb --username=root  ocr-database

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/ocr-database?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/ocr-database?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: createdb dropdb postgres migrateup migratedown sqlc test
