package main

import (
	"flag"
	"log"
	"os"

	"github.com/midorigreen/blobcat"
)

func run(bucket, prefix, ext string) error {
	s3 := blobcat.NewBlobS3(bucket, prefix, ext)
	return s3.ReadWrite(os.Stdout)
}

func main() {
	bucket := flag.String("b", "", "bucket name")
	prefix := flag.String("p", "", "prefix name")
	ext := flag.String("e", "gz", "extension name")
	flag.Parse()

	if err := run(*bucket, *prefix, *ext); err != nil {
		log.Fatal(err)
	}
}
