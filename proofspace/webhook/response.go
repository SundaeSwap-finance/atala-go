package webhook

func newWebhookResponse_Success(actionEventId string, serviceDid string, subscriberConnectDid string, issuedCredentials []WebhookCredentialValuesDTO) *WebhookResponse {
	return &WebhookResponse{
		Ok:                   true,
		Type:                 WebhookResponseSuccess,
		ActionEventId:        actionEventId,
		ProtocolVersion:      PROTOCOL_VERSION,
		ServiceDid:           serviceDid,
		SubscriberConnectDid: subscriberConnectDid,
		IssuedCredentials:    issuedCredentials,
	}
}
func newWebhookResponse_Failure(actionEventId string, errorMessage string) *WebhookResponse {
	return &WebhookResponse{
		Ok:              false,
		Type:            WebhookResponseFailure,
		ActionEventId:   actionEventId,
		ProtocolVersion: PROTOCOL_VERSION,
		Message:         errorMessage,
	}
}
