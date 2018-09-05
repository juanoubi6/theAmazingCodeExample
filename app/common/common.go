package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nats-io/go-nats"
	"github.com/streadway/amqp"
	"theAmazingCodeExample/app/config"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

var db *gorm.DB
var awsSession *session.Session
var rabbitMqConnection *amqp.Connection
var natsConnection *nats.Conn

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

func ConnectToRabbitMQ() {
	connection, err := amqp.Dial("amqp://" + config.GetConfig().RABBITMQ_USER + ":" + config.GetConfig().RABBITMQ_PASSWORD + "@" + config.GetConfig().RABBITMQ_HOST + ":" + config.GetConfig().RABBITMQ_PORT + "/")
	if err != nil {
		panic(err)
	}

	rabbitMqConnection = connection
}

func GetDatabase() *gorm.DB {
	return db
}

func GetRabbitMQChannel() *amqp.Channel {
	ch, err := rabbitMqConnection.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}

func CreateAWSSession() {

	creds := credentials.NewStaticCredentials(config.GetConfig().AWS_ACCESS_KEY_ID, config.GetConfig().AWS_SECRET_ACCESS_KEY, "")

	awsConfig := &aws.Config{
		Region: aws.String(config.GetConfig().AWS_REGION),
		Credentials: creds,
		Endpoint:aws.String(config.GetConfig().AWS_URL),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle:aws.Bool(true),

	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		panic(err)
	}

	awsSession = sess

	s3Client := s3.New(awsSession)

	// Create profile picture bucket
	cparams := &s3.CreateBucketInput{
		Bucket: aws.String(config.GetConfig().AWS_BUCKET_PROFILE_PICTURES),
	}
	_, err = s3Client.CreateBucket(cparams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				println("Bucket already owned")
			case s3.ErrCodeBucketAlreadyExists:
				println("Bucket already created")
			default:
				panic(err)
			}
		}
	}

}

func GetAWSSession() *session.Session {
	return awsSession
}

func ConnectToNats() {

	nc, err := nats.Connect(config.GetConfig().NATS_URL)
	if err != nil {
		panic(err)
	}

	natsConnection = nc

}

func GetNatsEncodedConnection() *nats.EncodedConn {

	c, err := nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	return c

}

func GetNatsConnection() *nats.Conn {
	return natsConnection
}
