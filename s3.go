package blobcat

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

const (
	region = "ap-northeast-1"

	// extention
	gzExt = "gz"
)

type s3ReadWrite struct {
	out io.WriteCloser
	m   sync.Mutex
}

func NewS3ReadWrite(out io.WriteCloser) *s3ReadWrite {
	return &s3ReadWrite{out: out}
}

func (w *s3ReadWrite) WriteAt(p []byte, off int64) (n int, err error) {
	w.m.Lock()
	defer w.m.Unlock()
	return w.out.Write(p)
}

func (w *s3ReadWrite) Close() error {
	w.m.Lock()
	defer w.m.Unlock()
	return w.out.Close()
}

type blobs3 struct {
	bucket string `require:"true"`
	prefix string `require:"true"`
	ext    string
}

func NewBlobS3(bucket, prefix, ext string) BlobReader {
	return &blobs3{
		bucket: bucket,
		prefix: prefix,
		ext:    ext,
	}
}

func (s *blobs3) ReadWrite(w io.Writer) error {
	var r = region
	sess := session.Must(session.NewSession(&aws.Config{Region: &r}))
	result, err := listObjects(s.bucket, s.prefix, sess)
	if err != nil {
		return err
	}
	downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.Concurrency = 1
	})
	for _, v := range result.Contents {
		rp, wp := io.Pipe()
		sw := NewS3ReadWrite(wp)

		input := &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    v.Key,
		}

		go func(input *s3.GetObjectInput, sw *s3ReadWrite) {
			defer sw.Close()
			err := download(input, sw, sess, downloader)
			if err != nil {
				log.Fatalf("download error : %v", err)
			}
		}(input, sw)

		writeExt(s.ext, rp, w)
	}

	return nil
}

func listObjects(bucket, prefix string, sess *session.Session) (*s3.ListObjectsV2Output, error) {
	svc := s3.New(sess)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
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

func writeExt(ext string, in io.Reader, out io.Writer) error {
	switch ext {
	case gzExt:
		gin, err := gzip.NewReader(in)
		if err != nil {
			return err
		}
		defer gin.Close()

		_, err = io.Copy(out, gin)
		if err != nil {
			return errors.Wrap(err, "gzip copy error")
		}
	default:
		return errors.New("not implements ext")
	}
	return nil
}
