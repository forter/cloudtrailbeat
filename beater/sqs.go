package beater

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func pullEvents(bt *Cloudtrailbeat) ([]CloudtrailRecord, error) {
	result, err := bt.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(bt.queueURL),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		return nil, err
	}
	bt.logger.Info("Number of SQS messages found", len(result.Messages))
	if len(result.Messages) == 0 {
		bt.logger.Info("No SQS Messages Found")
	}
	var toReturn []CloudtrailRecord
	for _, message := range result.Messages {
		body := message.Body
		handle := message.ReceiptHandle
		_, err := bt.sqs.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(bt.queueURL),
			ReceiptHandle: handle,
		})

		if err != nil {
			bt.logger.Error("Error acking message", err)
			return nil, err
		}
		var event events.SNSEntity
		err = json.Unmarshal([]byte(*body), &event)
		if err != nil {
			bt.logger.Error("Error decoding json", err)
			return nil, err
		}
		var s3Event events.S3Event
		err = json.Unmarshal([]byte(event.Message), &s3Event)
		buff := &aws.WriteAtBuffer{}
		bt.logger.Info("Attempting to download file", s3Event.Records)
		for _, r := range s3Event.Records {
			key, _ := url.Parse(r.S3.Object.Key)
			_, err = bt.downloader.Download(buff, &s3.GetObjectInput{
				Bucket: aws.String(r.S3.Bucket.Name),
				Key:    aws.String(key.Path),
			})
			if err != nil {
				bt.logger.Error("Error downloading file")
				return nil, err
			}
			var ctEvents CloudtrailFile
			gunzip, err := gzip.NewReader(bytes.NewBuffer(buff.Bytes()))
			if err != nil {
				bt.logger.Error("Couldnt create buffer", err)
				return nil, err
			}
			defer gunzip.Close()
			data, err := ioutil.ReadAll(gunzip)
			if err != nil {
				bt.logger.Error("Couldnt gunzip", err)
				return nil, err
			}
			_ = json.Unmarshal(data, &ctEvents)
			for _, record := range ctEvents.Records {
				toReturn = append(toReturn, record)
			}
		}
	}
	return toReturn, nil
}
