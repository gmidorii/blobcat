package blobcat

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

type S3 struct {
}

func (s *S3) Read(w io.WriterAt) error {
	// default setting only
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	obj := &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BLOB_BUCKET")),
		Key:    aws.String("sample.txt"),
	}

	err := download(obj, w, sess, downloader)
	if err != nil {
		return errors.Wrap(err, "download error")
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
