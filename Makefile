BIN=bc
CMD=cmd

# export env var or define below
#BLOB_BUCKET=
#BLOB_PREFIX=

build:
	cd ./$(CMD); go build -o $(BIN)

run: build
	./$(CMD)/$(BIN) -b $(BLOB_BUCKET) -p $(BLOB_PREFIX) -e gz
