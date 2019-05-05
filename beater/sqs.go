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
	logp.Info("Number of SQS messages found", len(result.Messages))
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
			logp.Err("Error acking message", err)
			return nil, err
		}
		var event events.SNSEntity
		err = json.Unmarshal([]byte(*body), &event)
		logp.Info("SNS Event", event)
		if err != nil {
			logp.Err("Error decoding json", err)
			return nil, err
		}
		var s3Event events.S3Event
		logp.Info("S3 Event", s3Event)
		err = json.Unmarshal([]byte(event.Message), &s3Event)
		buff := &aws.WriteAtBuffer{}
		logp.Info("Attempting to download file", s3Event.Records)
		for _, r := range s3Event.Records {
			logp.Info("File %s / %s", r.S3.Bucket.Name, r.S3.Object.Key)
			key, _ := url.Parse(r.S3.Object.Key)
			_, err = bt.downloader.Download(buff, &s3.GetObjectInput{
				Bucket: aws.String(r.S3.Bucket.Name),
				Key:    aws.String(key.Path),
			})
			if err != nil {
				logp.Err("Error downloading file")
				return nil, err
			}
			var ctEvents CloudtrailFile
			gunzip, err := gzip.NewReader(bytes.NewBuffer(buff.Bytes()))
			if err != nil {
				logp.Err("Couldnt create buffer", err)
				return nil, err
			}
			defer gunzip.Close()
			data, err := ioutil.ReadAll(gunzip)
			if err != nil {
				logp.Err("Couldnt gunzip", err)
				return nil, err
			}
			json.Unmarshal(data, &ctEvents)
			logp.Info("Should have events", ctEvents)
			for _, record := range ctEvents.Records {
				logp.Info("CT Record", record)
				toReturn = append(toReturn, record)
			}
		}
	}
	return toReturn, nil
}
