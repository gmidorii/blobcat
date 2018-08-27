BIN=bc
CMD=cmd

build:
	cd ./$(CMD); go build -o $(BIN)

run: build
	./$(CMD)/$(BIN)
