package webhook

import (
	"encoding/json"
	"net/http"
	"slices"
)

type ValidationError string

func (ce ValidationError) Error() string {
	return string(ce)
}

const MissingPublicServiceDid ValidationError = "Missing PublicServiceDid"
const MissingActionEventId ValidationError = "Missing ActionEventId"
const MissingSubscriberConnectDid ValidationError = "Missing SubscriberConnectDid"
const PublicServiceDidNotAllowed ValidationError = "PublicServiceDid not in allow list"

func decodeRequestParams(request *http.Request) (WebhookRequest, error) {
	var paramError error
	decoder := json.NewDecoder(request.Body)
	var requestBody WebhookRequest
	paramError = decoder.Decode(&requestBody)
	return requestBody, paramError
}

func ValidateRequestParams(request *http.Request, config *WebHookConfig) (*WebhookRequest, error) {
	requestBody, paramError := decodeRequestParams(request)
	if paramError != nil {
		return nil, paramError
	}
	if requestBody.PublicServiceDid == "" {
		paramError = MissingPublicServiceDid
	} else if !slices.Contains(config.AllowedServiceDids, requestBody.PublicServiceDid) {
		paramError = PublicServiceDidNotAllowed
	}
	if requestBody.ActionEventId == "" {
		paramError = MissingActionEventId
	}
	if requestBody.SubscriberConnectDid == "" {
		paramError = MissingSubscriberConnectDid
	}
	if paramError != nil {
		return nil, paramError
	}
	return &WebhookRequest{
		PublicServiceDid:     requestBody.PublicServiceDid,
		ActionEventId:        requestBody.ActionEventId,
		SubscriberConnectDid: requestBody.SubscriberConnectDid,
	}, paramError
}
