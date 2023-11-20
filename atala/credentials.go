package atala

import (
	"encoding/json"
	"fmt"
	"log"
)

const CRED_SERVICE_PATH = BASE_PATH + "issue-credentials/"

func (c *Client) CreateCredentialOffer(credOfferReq *CredentialOfferRequest) (*CredentialRecord, *ApiError, int, error) {
	body, err := json.Marshal(credOfferReq)
	if err != nil {
		log.Fatal("Error marshalling schema: ", err)
	}
	fmt.Println(credOfferReq)
	resp, credOffer, apiErr, err := PostRequest[CredentialRecord](c, CRED_SERVICE_PATH+"credential-offers", body)
	return credOffer, apiErr, resp.StatusCode, err
}
func (c *Client) AcceptCredentialOffer(holder_record_id string, subjectId string) (*CredentialRecord, *ApiError, int, error) {
	body := []byte(`{ "subjectId": "` + subjectId + `" }`)
	resp, credOffer, apiErr, err := PostRequest[CredentialRecord](c, CRED_SERVICE_PATH+"records/"+holder_record_id+"/accept-offer", body)
	return credOffer, apiErr, resp.StatusCode, err
}
func (c *Client) IssueCredential(issuer_record_id string) (*CredentialSchema, *ApiError, int, error) {
	resp, credOffer, apiErr, err := PostRequest[CredentialSchema](c, CRED_SERVICE_PATH+"records/"+issuer_record_id+"/issue-credential ", nil)
	return credOffer, apiErr, resp.StatusCode, err
}
func (c *Client) GetCredentialRecordsList() (*CredentialRecordList, *ApiError, int, error) {
	resp, credSchema, apiErr, err := GetRequest[CredentialRecordList](c, CRED_SERVICE_PATH+"records")
	return credSchema, apiErr, resp.StatusCode, err
}
func (c *Client) GetCredentialRecord(recordId string) (*CredentialRecord, *ApiError, int, error) {
	resp, credSchema, apiErr, err := GetRequest[CredentialRecord](c, CRED_SERVICE_PATH+"records/"+recordId)
	return credSchema, apiErr, resp.StatusCode, err
}
