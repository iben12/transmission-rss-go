build:
	go build -o ./bin/trss ./src

run:
	go run main.go

test:
	go test ./...