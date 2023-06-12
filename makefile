run:
	go run ./...
ls:
	sudo lsof -i:8081
.PHONEY: run