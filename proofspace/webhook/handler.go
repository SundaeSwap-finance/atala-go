package webhook

import (
	"fmt"
	"net/http"
)

const (
	PROTOCOL_VERSION = 1
)

/*
Handle a POST request to the ProofSpace "Interaction Webhook":
Verify request signature, Cleanse & Validate Request Parameters.
Return a WebhookResponse, success state issues new credentials; failure state with message;
see ProofSpace spec:
https://proofspace.atlassian.net/wiki/spaces/PSM/pages/2133786630/Integration+Webhooks+API+Overview
*/
func HandlePost(request *http.Request, config WebHookConfig) *WebhookResponse {
	// Verify Request Signature
	signatureVerificationErr := VerifyRequestSignature(request.Header, config.PublicKey)
	if signatureVerificationErr != nil {
		return newWebhookResponse_Failure(
			"", fmt.Sprintf("bad webhook request: %+v", signatureVerificationErr),
		)
	}
	// Cleanse & Validate Request Params
	webhookRequestBody, paramError := ValidateRequestParams(request, &config)
	if paramError != nil {
		errorResp := newWebhookResponse_Failure("", paramError.Error())
		return errorResp
	}

	// SUCCESS: Create credentials to be issued using config.IssueCredentials function
	var issuedCredentials []WebhookCredentialValuesDTO
	if config.IssueCredentials != nil {
		issuedCredentials = config.IssueCredentials(request, webhookRequestBody)
	}

	return newWebhookResponse_Success(
		webhookRequestBody.ActionEventId, webhookRequestBody.PublicServiceDid, webhookRequestBody.SubscriberConnectDid, issuedCredentials,
	)
}
