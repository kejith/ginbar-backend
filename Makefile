build:
	sqlc generate
	go build
run:
	go run main.go