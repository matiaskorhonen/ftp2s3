package main

import (
	"log"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/koofr/graval"
	"github.com/matiaskorhonen/ftp2s3/s3driver"
	"github.com/namsral/flag"
)

var host string
var port int
var username string
var password string
var serverName string

var awsRegion string
var awsAccessKeyID string
var awsSecretAccessKey string
var awsBucketName string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "host to bind to")
	flag.IntVar(&port, "port", 2121, "port to bind to")
	flag.StringVar(&username, "ftp-username", "ftp2s3", "FTP username")
	flag.StringVar(&password, "ftp-password", "ftp2s3", "FTP password")
	flag.StringVar(&serverName, "ftp-server-an,e", "FTP2S3", "FTP server name")

	flag.StringVar(&awsRegion, "aws-region", "us-east-1", "AWS region")
	flag.StringVar(&awsAccessKeyID, "aws-access-key-id", "", "AWS access key ID")
	flag.StringVar(&awsSecretAccessKey, "aws-secret-access-key", "", "AWS secret access key")
	flag.StringVar(&awsBucketName, "aws-bucket-name", "", "S3 bucket name")

	flag.String("config", "", "path to config file")

	flag.Parse()
}

func main() {
	// files := map[string]*s3driver.MemoryFile{
	// 	"/": &s3driver.MemoryFile{graval.NewDirItem(""), nil},
	// }

	credentialsProvider := aws.Creds(awsAccessKeyID, awsSecretAccessKey, "")

	factory := &s3driver.S3DriverFactory{
		Username:               username,
		Password:               password,
		AWSRegion:              awsRegion,
		AWSCredentialsProvider: credentialsProvider,
		AWSBucketName:          awsBucketName,
	}

	server := graval.NewFTPServer(&graval.FTPServerOpts{
		ServerName: serverName,
		Factory:    factory,
		Hostname:   host,
		Port:       port,
		PassiveOpts: &graval.PassiveOpts{
			ListenAddress: host,
			NatAddress:    host,
			PassivePorts: &graval.PassivePorts{
				Low:  42000,
				High: 45000,
			},
		},
	})

	log.Printf("FTP2S3 server listening on %s:%d", host, port)
	log.Printf("Access: ftp://%s:%s@%s:%d/", username, password, host, port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
