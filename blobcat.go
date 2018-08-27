package blobcat

import "io"

type BlobReader interface {
	Read(w io.WriterAt, bucket, key string) error
}
