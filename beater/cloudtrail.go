package beater

import (
	"encoding/json"
	
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/mitchellh/mapstructure"
)

type CloudtrailFile struct {
	Records []CloudtrailRecord
}

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-user-identity.html
type CloudtrailIdentity struct {
	Type                string                 `json:"type" mapstructure:"type"`
	UserName            string                 `json:"userName" mapstructure:"userName"`
	PrincipalID         string                 `json:"principalId" mapstructure:"principalId"`
	ARN                 string                 `json:"arn" mapstructure:"arn"`
	AccountID           string                 `json:"accountId" mapstructure:"accountId"`
	AccessKeyID         string                 `json:"accessKeyId" mapstructure:"accessKeyId"`
	SessionContext      map[string]interface{} `json:"sessionContext" mapstructure:"sessionContext"`
	InvokedBy           string                 `json:"invokedBy" mapstructure:"invokedBy"`
	SessionIssuer       map[string]interface{} `json:"sessionIssuer" mapstructure:"sessionIssuer"`
	WebIDFederationData map[string]interface{} `json:"webIdFederationData" mapstructure:"webIdFederationData"`
	IdentityProvider    string                 `json:"identityProvider" mapstructure:"identityProvider"`
}

func (cti *CloudtrailIdentity) ToCommonMap() (common.MapStr, error) {
	var result common.MapStr
	err := mapstructure.Decode(cti, &result)
	if err != nil {
		logp.L().Error("Error decoding Cloudtrail record", err)
		return nil, err
	}
	return result, nil
}

func (cti *CloudtrailIdentity) String() string {
	toReturn, err := json.MarshalIndent(cti, "", "    ")
	if err != nil {
		logp.L().Error("Error to json", err)
		return ""
	}
	return string(toReturn)
}

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-record-contents.html
type CloudtrailRecord struct {
	EventTime           string                   `json:"eventTime" mapstructure:"eventTime"`
	EventVersion        string                   `json:"eventVersion" mapstructure:"eventVersion"`
	UserIdentity        CloudtrailIdentity       `json:"userIdentity" mapstructure:"userIdentity"`
	EventSource         string                   `json:"eventSource" mapstructure:"eventSource"`
	EventName           string                   `json:"eventName" mapstructure:"eventName"`
	AWSRegion           string                   `json:"awsRegion" mapstructure:"awsRegion"`
	SourceIPAddress     string                   `json:"sourceIPAddress" mapstructure:"sourceIPAddress"`
	UserAgent           string                   `json:"userAgent" mapstructure:"userAgent"`
	ErrorCode           string                   `json:"errorCode" mapstructure:"errorCode"`
	ErrorMessage        string                   `json:"errorMessage" mapstructure:"errorMessage"`
	RequestParameters   map[string]interface{}   `json:"requestParameters" mapstructure:"requestParameters"`
	ResponseElements    string                   `json:"responseElements" mapstructure:"responseElements"`
	AdditionalEventData string                   `json:"additionalEventData" mapstructure:"additionalEventData"`
	RequestID           string                   `json:"requestID" mapstructure:"requestID"`
	EventID             string                   `json:"eventID" mapstructure:"eventID"`
	EventType           string                   `json:"eventType" mapstructure:"eventType"`
	ApiVersion          string                   `json:"apiVersion" mapstructure:"apiVersion"`
	ManagementEvent     string                   `json:"managementEvent" mapstructure:"managementEvent"`
	ReadOnly            string                   `json:"readOnly" mapstructure:"readOnly"`
	Resources           []map[string]interface{} `json:"resources" mapstructure:"resources"`
	RecipientAccountId  string                   `json:"recipientAccountId" mapstructure:"recipientAccountId"`
	ServiceEventDetails string                   `json:"serviceEventDetails" mapstructure:"serviceEventDetails"`
	SharedEventID       string                   `json:"sharedEventID" mapstructure:"sharedEventID"`
	VpcEndpointId       string                   `json:"vpcEndpointId" mapstructure:"vpcEndpointId"`
}

func (ctr *CloudtrailRecord) ToCommonMap() (common.MapStr, error) {
	var result common.MapStr
	err := mapstructure.Decode(ctr, &result)
	if err != nil {
		logp.L().Error("Error decoding Cloudtrail record", err)
		return nil, err
	}
	return result, nil
}

func (ctr *CloudtrailRecord) String() string{
	toReturn, err := json.MarshalIndent(ctr, "", "    ")
	if err != nil {
		logp.L().Error("Error to json", err)
		return ""
	}
	return string(toReturn)
}
