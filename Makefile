test:
	go test .
	
test-race:
	go test -race

race:
	go run -race .

build:
	go build .

exe: build
	./pipeline
	
run:
	go run .