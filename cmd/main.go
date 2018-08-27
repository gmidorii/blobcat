package main

import (
	"flag"
	"log"
	"os"

	"github.com/midorigreen/blobcat"
)

func run(bucket, key, ext string) error {
	s3 := blobcat.NewBlobS3(bucket, key, ext)
	return s3.ReadWrite(os.Stdout)
}

func main() {
	bucket := flag.String("b", "", "bucket name")
	key := flag.String("k", "", "key name")
	ext := flag.String("e", "gz", "extension name")
	flag.Parse()

	if err := run(*bucket, *key, *ext); err != nil {
		log.Fatal(err)
	}
}
