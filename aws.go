package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetEC2List(account *awsAuth) []*ec2.Reservation {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(account.Region),
		Credentials: credentials.NewStaticCredentials(account.AccessKey, account.SecretKey, ""),
	})
	if err != nil {
		log.Fatalf("AWS Session Error: %v", err)
	}

	svc := ec2.New(sess)
	result, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return nil
	}

	return result.Reservations
}

func FindEC2TagName(tags []*ec2.Tag) *string {
	for i, v := range tags {
		if *v.Key == "Name" {
			return tags[i].Value
		}
	}
	log.Println("Not Found Tag Name, Other Tags: ", tags)
	return nil
}

func AWSConsoleLink(awsAccount, sheetName *string) string {
	return "=HYPERLINK(\"https://" + *awsAccount + ".signin.aws.amazon.com/console\",\"" + *sheetName + "\")"
}