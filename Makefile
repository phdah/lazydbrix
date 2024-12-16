BIN := bin/lazydbrix
SOURCE := cmd/lazydbrix/main.go

build:
	go build -o $(BIN) $(SOURCE)

run: build
	$(BIN)

clean:
	rm $(BIN) || echo "No bin exists"
