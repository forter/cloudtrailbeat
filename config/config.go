// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration `config:"period"`
	SQSQueueName string        `config:"sqs_queue_name"`
	AccountID    string        `config:"account_id"`
}

var DefaultConfig = Config{
	Period:       1 * time.Second,
	SQSQueueName: "",
	AccountID:    "",
}
