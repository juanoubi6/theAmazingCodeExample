package amazonS3

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"hash/adler32"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/config"
	"time"
)

func UploadImageToS3(header *multipart.FileHeader) (string, string, error) {

	if err := checkImageType(header); err != nil {
		return "", "", err
	}

	s3Key, url, err := uploadPicture(common.GetAWSSession(), header)
	if err != nil {
		return "", "", err
	}

	return s3Key, url, nil

}

func checkImageType(header *multipart.FileHeader) error {
	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	ctype, err := getContentType(file)
	if err != nil {
		return err
	}
	validCTypes := map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
	}
	if _, ok := validCTypes[ctype]; !ok {
		return errors.New(header.Filename + ": invalid image type")
	}

	return nil
}

func getContentType(file multipart.File) (string, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}

	return http.DetectContentType(buff), nil
}

func uploadPicture(awsSession *session.Session, header *multipart.FileHeader) (string, string, error) {

	file, err := header.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	ctype, err := getContentType(file)
	if err != nil {
		return "", "", err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		return "", "", err
	}

	uploader := s3manager.NewUploader(awsSession)
	key := aws.String(getS3FileKey(buf.Bytes()) + "." + ctype[6:])
	uploaded, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(config.GetConfig().AWS_BUCKET),
		Key:         key,
		Body:        file,
		ContentType: aws.String(ctype),
	})
	if err != nil {
		return "", "", err
	}

	return *key, uploaded.Location, nil

}

func getS3FileKey(file []byte) string {
	date := []byte(time.Now().String())
	cs := adler32.Checksum(append(file, date...))
	return strconv.Itoa(int(cs))
}
