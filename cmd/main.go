package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/midorigreen/blobcat"
)

const (
	bufSize = 1000
	gzExt   = "gz"
)

func run(bucket, key, ext string) error {
	s3 := blobcat.NewBlobS3()
	buf := make([]byte, bufSize)
	bufAt := aws.NewWriteAtBuffer(buf)
	err := s3.Read(bufAt, bucket, key)
	if err != nil {
		return err
	}

	switch ext {
	case gzExt:
		rb := bytes.NewBuffer(bufAt.Bytes())
		gr, err := gzip.NewReader(rb)
		if err != nil {
			return err
		}
		defer gr.Close()

		io.Copy(os.Stdout, gr)
	default:
		fmt.Fprint(os.Stdout, string(bufAt.Bytes()))
	}
	return nil
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
