package blobcat

import "io"

type BlobReader interface {
	ReadWrite(w io.Writer) error
}
