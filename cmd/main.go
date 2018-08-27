package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/midorigreen/blobcat"
)

func run(bucket, key string) error {
	s3 := blobcat.NewBlobS3()
	buf := make([]byte, 1000)
	bufAt := aws.NewWriteAtBuffer(buf)
	err := s3.Read(bufAt, bucket, key)
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, string(bufAt.Bytes()))
	return nil
}

func main() {
	bucket := flag.String("b", "", "bucket name")
	key := flag.String("k", "", "key name")
	flag.Parse()

	if err := run(*bucket, *key); err != nil {
		log.Fatal(err)
	}
}
