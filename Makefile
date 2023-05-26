postgres:
	sudo docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createbd: 
	sudo docker exec -it postgres12 createdb --username=root --owner=root course
dropdb:
	sudo docker exec -it postgres12 dropdb  course

migrateup: 
	migrate -path migration/ -database "postgresql://root:secret@localhost:5432/course?sslmode=disable" -verbose up

migratedown:
	migrate -path migration/ -database "postgresql://root:secret@localhost:5432/course?sslmode=disable" -verbose down
.PHONY: postgres createbd dropdb migrateup migratedown