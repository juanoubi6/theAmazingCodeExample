package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mime/multipart"
	"sync"
	"theAmazingCodeExample/app/config"
)

var db *gorm.DB
var awsSession *session.Session

func ConnectToDatabase() {
	var err error
	dbname := config.GetConfig().DB_NAME
	dbhost := config.GetConfig().DB_HOST
	dbport := config.GetConfig().DB_PORT
	dbuser := config.GetConfig().DB_USERNAME
	dbpass := config.GetConfig().DB_PASSWORD

	db, err = gorm.Open("mysql", dbuser+":"+dbpass+"@"+"tcp("+dbhost+":"+dbport+")"+"/"+dbname+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

}

func GetDatabase() *gorm.DB {
	return db
}

func CreateAWSSession() {
	creds := credentials.NewStaticCredentials(config.GetConfig().AWS_ACCESS_KEY_ID, config.GetConfig().AWS_SECRET_ACCESS_KEY, "")

	sess, err := session.NewSession(&aws.Config{Region: aws.String(config.GetConfig().AWS_REGION),
		Credentials: creds})
	if err != nil {
		println(err.Error())
	}

	awsSession = sess
}

func GetAWSSession() *session.Session {
	return awsSession
}

/////Worker pool//////
type Task interface {
	run(wg *sync.WaitGroup)
}

type UploadImageTask struct {
	FileHeader *multipart.FileHeader
	UserID     uint
	Err        error
	Function   func(*multipart.FileHeader, uint) error
}

type Pool struct {
	Tasks        []*Task
	Concurrency  int
	TasksChannel chan *Task
	Wg           sync.WaitGroup
}

func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:        tasks,
		Concurrency:  concurrency,
		TasksChannel: make(chan *Task),
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
