BIN=bc
CMD=cmd

# export env var or define below
#BLOB_BUCKET=
#BLOB_KEY=

build:
	cd ./$(CMD); go build -o $(BIN)

run: build
	./$(CMD)/$(BIN) -b $(BLOB_BUCKET) -k $(BLOB_KEY)
