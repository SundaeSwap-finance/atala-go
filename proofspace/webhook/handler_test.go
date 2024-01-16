package webhook_test

import (
	. "atala-go/proofspace/webhook"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_HandlePost_Success(t *testing.T) {
	actionEventId := "actionEventId_asdfg"
	serviceDid := "serviceDid_123345"
	subscriberConnectDid := "subscriberConnectDid_54321"
	aMeaningfulHash := "some_meaninfulHash_12345asdfg"
	wantCreds := newCredToIssue(aMeaningfulHash)
	receivedCredentials := newCredToReceive()
	webhookRequestBody := WebhookRequest{
		ActionEventId:        actionEventId,
		PublicServiceDid:     serviceDid,
		SubscriberConnectDid: subscriberConnectDid,
		ReceivedCredentials:  receivedCredentials,
	}
	req, publicKey := newMockRequest(t, webhookRequestBody)
	webhookConfig := WebHookConfig{
		IssueCredentials:   create_IssueCredentialsFn(wantCreds),
		AllowedServiceDids: []string{serviceDid},
		PublicKey:          publicKey,
	}
	gotResponse := HandlePost(req, webhookConfig)
	wantResponse := &WebhookResponse{
		Ok:                   true,
		Type:                 "success",
		ServiceDid:           serviceDid,
		SubscriberConnectDid: subscriberConnectDid,
		ActionEventId:        actionEventId,
		ProtocolVersion:      PROTOCOL_VERSION,
		Message:              "",
		IssuedCredentials:    wantCreds,
		RevokedCredentials:   nil,
	}
	assert.Equal(t, wantResponse, gotResponse)
}

func Test_HandlePost_Failure_missingFields(t *testing.T) {
	actionEventId := "actionEventId_asdfg"
	serviceDid := "serviceDid_123345"
	subscriberConnectDid := "subscriberConnectDid_54321"
	// aMeaningfulHash := "some_meaningfulHash_12345asdfg"
	wantResponse := &WebhookResponse{
		Ok:                   false,
		Type:                 "failure",
		ServiceDid:           "",
		SubscriberConnectDid: "",
		ActionEventId:        "",
		ProtocolVersion:      PROTOCOL_VERSION,
		Message:              "ERROR: missing ActionEventId",
		IssuedCredentials:    nil,
		RevokedCredentials:   nil,
	}
	issueCredentialsFn_ReturnCreds := newCredToIssue("this cred will not be issued")
	receivedCredentials := newCredToReceive()
	webhookRequestBody := WebhookRequest{
		ActionEventId:        actionEventId,
		PublicServiceDid:     serviceDid,
		SubscriberConnectDid: subscriberConnectDid,
		ReceivedCredentials:  receivedCredentials,
	}
	testFunc := func(fieldToClear string, body WebhookRequest, wantMessage string, serviceDid_Opt ...string) {
		_serviceDid := serviceDid
		if len(serviceDid_Opt) > 0 {
			_serviceDid = serviceDid_Opt[0]
		}
		if fieldToClear != "" {
			reflect.ValueOf(&body).Elem().FieldByName(fieldToClear).Set(reflect.ValueOf(""))
		}
		req, pubKey := newMockRequest(t, body)
		webhookConfig := WebHookConfig{
			IssueCredentials:   create_IssueCredentialsFn(issueCredentialsFn_ReturnCreds),
			AllowedServiceDids: []string{_serviceDid},
			PublicKey:          pubKey,
		}
		// t.Logf("%v\n%++v\n%++v\n", fieldToClear, req, webhookConfig)
		gotResponse := HandlePost(req, webhookConfig)
		wantResponse.Message = wantMessage
		assert.Equal(t, wantResponse, gotResponse)
	}
	testFunc("PublicServiceDid", webhookRequestBody, MissingPublicServiceDid.Error())
	testFunc("ActionEventId", webhookRequestBody, MissingActionEventId.Error())
	testFunc("SubscriberConnectDid", webhookRequestBody, MissingSubscriberConnectDid.Error())
	testFunc("", webhookRequestBody, PublicServiceDidNotAllowed.Error(), "the_only_allowed_serviceDid")
}

/* TEST HELPERS */

/* Given a WebhookRequest object, creates a mock HTTP POST request with a valid JWT auth token in header */
func newMockRequest(t *testing.T, webhookRequestBody WebhookRequest) (*http.Request, string) {
	path := "[any]"
	u := &url.URL{Path: path}
	// Prepare request body
	b, err := json.Marshal(webhookRequestBody)
	if err != nil {
		t.Errorf("error marshalling webhookRequestBody: %++v", err)
	}
	var buf io.ReadWriter
	if b != nil {
		buf = bytes.NewBuffer([]byte(b))
	}
	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		t.Errorf("error creating new request object: %++v", err)
	}
	// Generate test keys and jwt token
	rsakey := newTestKey(t)
	publicPEM := encodePublicPEM(rsakey)
	tokenString := newTokenString(t, rsakey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(AuthHeaderName, "Bearer "+tokenString)
	return req, string(publicPEM)
}

/* Creates a Config.IssueCredentials function which returns the passed credentials */
func create_IssueCredentialsFn(creds []WebhookCredentialValuesDTO) IssueCredenialsFunction {
	return func(r *http.Request, params *WebhookRequest) []WebhookCredentialValuesDTO {
		return creds
	}
}

/* Creates a WebhookCredentialValuesDTO to be issued by the IssueCredentials func */
func newCredToIssue(aMeaningfulHash string) []WebhookCredentialValuesDTO {
	return []WebhookCredentialValuesDTO{
		NewCredential(
			"some_credential_to_issue:1:AB:101:tag",
			[]WebhookCredentialField{
				{
					Name:  "Credential Issue Date",
					Value: fmt.Sprint(time.Now().UTC().UnixMilli()),
				},
				{
					Name:  "The Hash",
					Value: aMeaningfulHash,
				},
			},
		),
	}
}

/* Creates a dummy WebhookCredentialValuesDTO to be received in the request body as input to the webhook */
func newCredToReceive() []WebhookCredentialValuesDTO {
	return []WebhookCredentialValuesDTO{
		NewCredential(
			"Some_Received_Cred",
			[]WebhookCredentialField{
				{
					Name:  "Credential Issue Date",
					Value: fmt.Sprint(time.Now().UTC().UnixMilli()),
				},
				{
					Name:  "A Field",
					Value: "some value",
				},
			},
		),
	}
}
