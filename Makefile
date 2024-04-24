test:
	go test -race

race:
	go run -race .

run:
	go build . | ./pipeline