package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/midorigreen/blobcat"
)

func run() error {
	s3 := &blobcat.S3{}
	buf := make([]byte, 1000)
	bufAt := aws.NewWriteAtBuffer(buf)
	err := s3.Read(bufAt)
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, string(bufAt.Bytes()))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
