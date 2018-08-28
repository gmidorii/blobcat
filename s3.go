package blobcat

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	var r = region
	sess := session.Must(session.NewSession(&aws.Config{Region: &r}))
	result, err := listObjects(s.bucket, s.key, sess)
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess)
	for _, v := range result.Contents {
		input := &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    v.Key,
		}

		buf := make([]byte, bufSize)
		bufAt := aws.NewWriteAtBuffer(buf)
		err := download(input, bufAt, sess, downloader)
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

			io.Copy(w, gr)
		default:
			fmt.Fprint(w, string(bufAt.Bytes()))
		}
	}

	return nil
}

func listObjects(bucket, key string, sess *session.Session) (*s3.ListObjectsV2Output, error) {
	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, fmt.Errorf("%v: %v", s3.ErrCodeNoSuchBucket, aerr)
			default:
				return nil, aerr
			}
		} else {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func download(obj *s3.GetObjectInput, w io.WriterAt, sess *session.Session, downloader *s3manager.Downloader) error {
	_, err := downloader.Download(w, obj)
	if err != nil {
		return errors.Wrap(err, "failed downloader download")
	}
	return nil
}
