package utils

import (
	"bytes"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/hiyali/logli"
)

const (
	S3_REGION = "us-east-2"
	S3_BUCKET = "eiprice.delivery"
	S3_ACL    = "public-read"
)

type S3Handler struct {
	Session *session.Session
	Bucket  string
}

func UploadCSV(key string, filename string) {

	AccessKeyID := "AKIAIEFOZ7BVF4ZSBRVA"
	SecretAccessKey := "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"

	file, err := os.Open(filename)
	if err != nil {
		log.FatalF("os.Open - filename: %v, err: %v", filename, err)
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			AccessKeyID,
			SecretAccessKey,
			"",
		),
	})
	if err != nil {
		log.FatalF("session.NewSession - filename: %v, err: %v", filename, err)
	}

	handler := S3Handler{
		Session: sess,
		Bucket:  S3_BUCKET,
	}

	err = handler.UploadFile(key, filename)
	if err != nil {
		log.FatalF("UploadFile - filename: %v, err: %v", filename, err)
	}

	log.Info(`##### url s3: https://s3.us-east-2.amazonaws.com/eiprice.delivery/` + key)
	log.Info("Upload to s3 - success")
}

func (h S3Handler) UploadFile(key string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.FatalF("os.Open - filename: %s, err: %v", filename, err)
	}
	defer file.Close()

	_, err = s3.New(h.Session).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(h.Bucket),
		Key:                aws.String(key),
		ACL:                aws.String(S3_ACL),
		Body:               file,
		ContentDisposition: aws.String("attachment"),
	})

	return err
}

func (h S3Handler) ReadFile(key string) (string, error) {
	results, err := s3.New(h.Session).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(h.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	defer results.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, results.Body); err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}
