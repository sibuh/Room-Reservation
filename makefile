run:
	go run ./...
sqlc:
	sqlc generate 
.PHONEY: run