package blobcat

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

const region = "ap-northeast-1"

type blobs3 struct {
}

func NewBlobS3() BlobReader {
	return &blobs3{}
}

func (s *blobs3) Read(w io.WriterAt, bucket, key string) error {
	// default setting only
	var r = region
	sess := session.Must(session.NewSession(&aws.Config{Region: &r}))
	downloader := s3manager.NewDownloader(sess)
	obj := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
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
