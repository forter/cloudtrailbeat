package beater

import (
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/mitchellh/mapstructure"
)

type CloudtrailFile struct {
	Records []CloudtrailRecord
}

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-user-identity.html
type CloudtrailIdentity struct {
	Type                string                 `json:"type"`
	UserName            string                 `json:"userName"`
	PrincipalID         string                 `json:"principalId"`
	ARN                 string                 `json:"arn"`
	AccountID           string                 `json:"accountId"`
	AccessKeyID         string                 `json:"accessKeyId"`
	SessionContext      map[string]interface{} `json:"sessionContext"`
	InvokedBy           string                 `json:"invokedBy"`
	SessionIssuer       map[string]interface{} `json:"sessionIssuer"`
	WebIDFederationData map[string]interface{} `json:"webIdFederationData"`
	IdentityProvider    string                 `json:"identityProvider"`
}

func (cti *CloudtrailIdentity) ToCommonMap() (common.MapStr, error) {
	var result common.MapStr
	err := mapstructure.Decode(cti, &result)
	if err != nil {
		logp.Err("Error decoding Cloudtrail record", err)
		return nil, err
	}
	return result, nil
}

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-record-contents.html
type CloudtrailRecord struct {
	EventTime           string                   `json:"eventTime"`
	EventVersion        string                   `json:"eventVersion"`
	UserIdentity        CloudtrailIdentity       `json:"userIdentity"`
	EventSource         string                   `json:"eventSource"`
	EventName           string                   `json:"eventName"`
	AWSRegion           string                   `json:"awsRegion"`
	SourceIPAddress     string                   `json:"sourceIPAddress"`
	UserAgent           string                   `json:"userAgent"`
	ErrorCode           string                   `json:"errorCode"`
	ErrorMessage        string                   `json:"errorMessage"`
	RequestParameters   map[string]interface{}   `json:"requestParameters"`
	ResponseElements    string                   `json:"responseElements"`
	AdditionalEventData string                   `json:"additionalEventData"`
	RequestID           string                   `json:"requestID"`
	EventID             string                   `json:"eventID"`
	EventType           string                   `json:"eventType"`
	ApiVersion          string                   `json:"apiVersion"`
	ManagementEvent     string                   `json:"managementEvent"`
	ReadOnly            string                   `json:"readOnly"`
	Resources           []map[string]interface{} `json:"resources"`
	RecipientAccountId  string                   `json:"recipientAccountId"`
	ServiceEventDetails string                   `json:"serviceEventDetails"`
	SharedEventID       string                   `json:"sharedEventID"`
	VpcEndpointId       string                   `json:"vpcEndpointId"`
}

func (ctr *CloudtrailRecord) ToCommonMap() (common.MapStr, error) {
	var result common.MapStr
	err := mapstructure.Decode(ctr, &result)
	if err != nil {
		logp.Err("Error decoding Cloudtrail record", err)
		return nil, err
	}
	return result, nil
}
