build:
	go build -o ./bin/trss ./main.go

run:
	go run main.go

test:
	go test -v -cover -parallel 2 ./trss...