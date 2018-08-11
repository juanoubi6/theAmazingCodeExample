package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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


