package blobcat

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

const (
	region  = "ap-northeast-1"
	bufSize = 1000
	gzExt   = "gz"
)

type blobs3 struct {
	bucket string `require:"true"`
	key    string `require:"true"`
	ext    string
}

func NewBlobS3(bucket, key, ext string) BlobReader {
	return &blobs3{
		bucket: bucket,
		key:    key,
		ext:    ext,
	}
}

func (s *blobs3) ReadWrite(w io.Writer) error {
	// default setting only
	buf := make([]byte, bufSize)
	bufAt := aws.NewWriteAtBuffer(buf)

	var r = region
	sess := session.Must(session.NewSession(&aws.Config{Region: &r}))
	downloader := s3manager.NewDownloader(sess)
	obj := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key),
	}

	err := download(obj, bufAt, sess, downloader)
	if err != nil {
		return errors.Wrap(err, "download error")
	}

	switch s.ext {
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

func download(obj *s3.GetObjectInput, w io.WriterAt, sess *session.Session, downloader *s3manager.Downloader) error {
	_, err := downloader.Download(w, obj)
	if err != nil {
		return errors.Wrap(err, "failed downloader download")
	}
	return nil
}
