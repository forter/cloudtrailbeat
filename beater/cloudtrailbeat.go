package beater

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/forter/cloudtrailbeat/config"
)

// Cloudtrailbeat configuration.
type Cloudtrailbeat struct {
	done       chan struct{}
	config     config.Config
	client     beat.Client
	sqs        *sqs.SQS
	queueURL   string
	downloader *s3manager.Downloader
}

// New creates an instance of cloudtrailbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	sess := session.Must(session.NewSession())
	svc := sqs.New(sess)
	queueURLResp, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(c.SQSQueueName),
	})
	if err != nil {
		logp.Err("Could not get Queue Name")
		return nil, err
	}
	s3Svc := s3.New(sess)
	downloader := s3manager.NewDownloaderWithClient(s3Svc)
	bt := &Cloudtrailbeat{
		done:       make(chan struct{}),
		config:     c,
		sqs:        svc,
		queueURL:   queueURLResp.String(),
		downloader: downloader,
	}
	return bt, nil
}

// Run starts cloudtrailbeat.
func (bt *Cloudtrailbeat) Run(b *beat.Beat) error {
	logp.Info("cloudtrailbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		// poll sqs queue
		events, err := pullEvents(bt)
		if err != nil {
			logp.Err("Failed to pull events from SQS", err)
		}
		fields := common.MapStr{}
		fields["type"] = b.Info.Name
		for _, e := range events {
			values, err := e.ToCommonMap()
			if err != nil {
				logp.Err("Shittt")
			}
			event := beat.Event{
				Timestamp: time.Now(),
				Fields:    values,
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
		}
	}
}

// Stop stops cloudtrailbeat.
func (bt *Cloudtrailbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
