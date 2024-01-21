# ProofSpace WebHook Utilities

Provides utilities for ProofSpace integration via WebHooks.

## Prerequisite: Set Up Your ProofSpace Dashboard

In the [ProofSpace Dashboard](https://sso.proofspace.id/), create the neccessary service, schemas, credentials, and interactions.

## Usage

The [webhook](webhook/) module can be used in 2 ways:

1. Calling the ```HandlePost(r *http.Request, config WebHookConfig)``` function in [handler.go](webhook/handler.go)
2. Calling the various module functions (```VerifyRequestSignature```, ```ValidateRequestParams```, etc.) individually.

```HandlePost``` performs the following processing on the request and response:  

* JWT verification via [VerifyRequestSignature](webhook/verify_request_sig.go)
* Cleanse &amp; validate ```WebhookRequest``` Parameters via [ValidateRequestParams](webhook/validation.go)
* Determing whether or not to issue a new credential via config.IssueCredentials [IssueCredentialsFunction](webhook/types.go)

### Usage via ```HandlePost```

When a new Request is received, e.g.:
```
func OnNewRequest(writer http.ResponseWriter, r *http.Request)
```
1. First, define the 3 values needed for configuration:  
    ```SERVICE_DID``` &amp; ```PUBLIC_KEY```:
    ```
    const (
        SERVICE_DID = "some_service_did_12345"
        PUBLIC_KEY  = `-----BEGIN PUBLIC ... full contents of public key file ... `
    )
    ```
    ```issueCredentials```:
    ```
    func issueCredentials(
        request *http.Request, params *webhook.WebhookRequestParams
    ) []webhook.WebhookCredentialValuesDTO {
        /* Perform logic to determine whether to create a credential */
        if shouldIssue {
            // Create a credential to be issued
            return []webhook.WebhookCredentialValuesDTO{
                webhook.NewCredential(
                    "SOME_CREDENTIAL_ID",
                    []webhook.WebhookCredentialField{
                        {
                            Name:  "Credential Issue Date",
                            Value: fmt.Sprint(time.Now().UTC().UnixMilli()),
                        },
                    },
                ),
            }
        }
        // Otherwise, no credentials will be issued
        return []webhook.WebhookCredentialValuesDTO{}
    }
    ```
2. Next, create a ```WebHookConfig``` object with the above values:
    ```
    webhookConfig := webhook.WebHookConfig{
        AllowedServiceDids: []string{serviceDid},
        IssueCredentials:   issueCredentials,
        PublicKey:          PUBLIC_KEY,
    }
    ```
    Note the type definition:
    ```
    type WebHookConfig struct {
        AllowedServiceDids []string
        PublicKey          string
        // func called after successful sig verification & param validation
        IssueCredentials   IssueCredentialsFunction
    }
    ```
3. Lastly, pass the request and config into ```HandlePost```, &amp; handle the JSON Response:
    ```
    resp := webhook.HandlePost(r, webhookConfig)
	respJson, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(writer, "json error")
		return
	}
	fmt.Fprintf(writer, "%s", respJson)
    ```

## Generate QR Codes

A common way to initiate interactions is via QR Code. These QR Codes and are scanned by the ProofSpace wallet to launch an interaction. When the user confirms the interaction, the WebHook receives a Request.

ProofSpace Wiki [How-to integrate QR/Web-link to trigger Interaction](https://proofspace.atlassian.net/wiki/spaces/PSM/pages/2196930561/How-to+integrate+QR+Web-link+to+trigger+Interaction)

## Complete Samples

Create a function app for your webhook
```coming soon```


## Run Tests

```
cd proofspace/webhook
go test -v
```