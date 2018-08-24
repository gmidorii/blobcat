package main

import (
	"log"
	"os"

	"github.com/midorigreen/blobcat"
)

func run() error {
	s3 := &blobcat.S3{}
	return s3.Read(os.Stdout)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
