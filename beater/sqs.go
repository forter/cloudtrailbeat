package beater

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/elastic/beats/libbeat/logp"
)

func pullEvents(bt *Cloudtrailbeat) ([]CloudtrailRecord, error) {
	result, err := bt.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(bt.config.SQSUrl),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		return nil, err
	}
	if len(result.Messages) == 0 {
		logp.Info("No SQS Messages Found")
	}
	toReturn := []CloudtrailRecord{}
	for _, message := range result.Messages {
		body := message.Body
		handle := message.ReceiptHandle
		_, err := bt.sqs.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &bt.config.SQSUrl,
			ReceiptHandle: handle,
		})

		if err != nil {
			fmt.Println("Delete Error", err)
			return nil, err
		}
		var event events.SNSEventRecord
		err = json.Unmarshal([]byte(*body), &event)
		if err != nil {
			fmt.Println("Error decoding json")
			return nil, err
		}
		var s3Event events.S3EventRecord
		err = json.Unmarshal([]byte(event.SNS.Message), &s3Event)
		buff := &aws.WriteAtBuffer{}
		_, err = bt.downloader.Download(buff, &s3.GetObjectInput{
			Bucket: aws.String(s3Event.S3.Bucket.Name),
			Key:    aws.String(s3Event.S3.Object.Key),
		})
		var ctEvents CloudtrailFile
		json.Unmarshal(buff.Bytes(), ctEvents)
		for _, record := range ctEvents.Records {
			toReturn = append(toReturn, record)
		}
	}
	return toReturn, nil
}
