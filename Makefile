BIN := bin/lazydbrix

build:
	go build -o $(BIN) cmd/main.go

clean:
	rm $(BIN) || echo "No bin exists"
