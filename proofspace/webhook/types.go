package webhook

import "net/http"

type IssueCredenialsFunction = func(request *http.Request, params *WebhookRequest) []WebhookCredentialValuesDTO

type WebHookConfig struct {
	// func called after successful sig verification & param validation
	IssueCredentials   IssueCredenialsFunction
	AllowedServiceDids []string
	PublicKey          string
}

type WebhookRequest struct {
	PublicServiceDid     string                       `json:"publicServiceDid"`     // your service DID
	SubscriberConnectDid string                       `json:"subscriberConnectDid"` // connect DID of the user
	ActionId             string                       `json:"actionId"`             // interaction id from dashboard
	ActionInstanceId     string                       `json:"actionInstanceId"`     // interaction instance id
	ActionEventId        string                       `json:"actionEventId"`        // id of the event  (unique for the service)
	ActionParams         []WebhookCredentialField     `json:"actionParams"`         // array of parameters, configured in action
	ReceivedCredentials  []WebhookCredentialValuesDTO `json:"receivedCredentials"`  // array of required credentials
}

type WebhookResponse struct {
	Ok   bool                `json:"ok"`
	Type WebhookResponseType `json:"type"` // failure | success

	ServiceDid           string `json:"serviceDid"`           // public service DID, should be the same as in request
	SubscriberConnectDid string `json:"subscriberConnectDid"` // connect DID, should be the same as in request
	ActionEventId        string `json:"actionEventId"`        // event
	ProtocolVersion      int    `json:"protocolVersion"`      //
	Message              string `json:"message"`              // error message

	// array of issued credentials in case you want to issue them immediately. See the structure at https://proofspace.atlassian.net/wiki/spaces/PSM/pages/2133786630/Integration+Webhooks+API+Overview
	IssuedCredentials []WebhookCredentialValuesDTO `json:"issuedCredentials"`
	// can be kept as an empty array, not used
	RevokedCredentials []WebhookCredentialValuesDTO `json:"revokedCredentials"`
}

type WebhookResponseType string

const (
	WebhookResponseFailure WebhookResponseType = "failure"
	WebhookResponseSuccess WebhookResponseType = "success"
)

type WebhookCredentialValuesDTO struct {
	CredentialId string                   `json:"credentialId"`
	SchemaId     string                   `json:"schemaId"`
	Fields       []WebhookCredentialField `json:"fields"`       // credentials from phone
	UtcIssuedAt  int                      `json:"utcIssuedAt"`  // utc   time in milliseconds
	Revoked      bool                     `json:"revoked"`      // is this credential revoked.  (optional, not used)
	UtcRevokedAt int                      `json:"utcRevokedAt"` // utc   time in milliseconds
}
type WebhookCredentialField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CredentialConfig struct {
	SchemaId     string
	CredentialId string
}

/*
FROM: https://proofspace.atlassian.net/wiki/spaces/PSM/pages/2133786630/Integration+Webhooks+API+Overview

* Note: the fields values in WebhookCredentialValuesDTO[ ] should all be string.
According to the attribute type, which is set in the credential schema, values should be set to:
Text —  text itself
Date — string with the number of seconds since epoch.
Number — number as string.  If this number is float, than '.'  used as decimal point   "222.2222"
ImageURL -- url as string.  // NOT USED YET
Enumeration -- numeric enumeration value as string, i.e.  "0", "1",..etc
Actual value types will be casted on credential issue according to the credential schema attribute types.
*/
