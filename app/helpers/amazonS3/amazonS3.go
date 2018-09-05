package amazonS3

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"hash/adler32"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"theAmazingCodeExample/app/common"
	"theAmazingCodeExample/app/config"
	"theAmazingCodeExample/app/models"
	"time"
)

func DeletePictureFromS3(photoData models.ProfilePicture, bucketName string) error {

	if config.GetConfig().AWS_SECRET_ACCESS_KEY != "" {

		svc := s3.New(common.GetAWSSession())

		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(photoData.S3Key),
		})
		if err != nil {
			return err
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(photoData.S3Key),
		})
		if err != nil {
			return err
		}

	}

	err := common.GetDatabase().Delete(&photoData).Error
	if err != nil {
		return err
	}

	return nil
}

func UploadImageToS3(header *multipart.FileHeader, bucketName string) (string, string, error) {

	if err := checkImageType(header); err != nil {
		return "", "", err
	}

	s3Key, url, err := uploadPicture(common.GetAWSSession(), header, bucketName)
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

func uploadPicture(awsSession *session.Session, header *multipart.FileHeader, bucketName string) (string, string, error) {

	if config.GetConfig().AWS_SECRET_ACCESS_KEY != "" {

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
			Bucket:      aws.String(bucketName),
			Key:         key,
			Body:        file,
			ContentType: aws.String(ctype),
		})
		if err != nil {
			return "", "", err
		}

		return *key, uploaded.Location, nil

	} else {
		return "falseS3key", "falseURL", nil
	}

}

func getS3FileKey(file []byte) string {
	date := []byte(time.Now().String())
	cs := adler32.Checksum(append(file, date...))
	return strconv.Itoa(int(cs))
}

///////////////////Worker pool implementation///////////////////
type UploadImageTask struct {
	FileHeader *multipart.FileHeader
	UserID     uint
	Err        error
	Function   func(*multipart.FileHeader, uint) error
}

type Pool struct {
	Tasks        []*UploadImageTask
	Concurrency  int
	TasksChannel chan *UploadImageTask
	Wg           sync.WaitGroup
}

func NewPool(tasks []*UploadImageTask, concurrency int) *Pool {
	return &Pool{
		Tasks:        tasks,
		Concurrency:  concurrency,
		TasksChannel: make(chan *UploadImageTask),
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.Concurrency; i++ {
		go p.work()
	}

	p.Wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.TasksChannel <- task
	}

	close(p.TasksChannel)

	p.Wg.Wait()
}

func (p *Pool) work() {
	for task := range p.TasksChannel {
		task.run(&p.Wg)
	}
}

func (t *UploadImageTask) run(wg *sync.WaitGroup) {
	t.Err = t.Function(t.FileHeader, t.UserID)
	wg.Done()
}
