run:
	go run ./...
ls:
	sudo lsof -i:8081
sqlc:
	sqlc generate 
.PHONEY: run