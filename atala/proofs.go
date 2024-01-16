package atala

import (
	"encoding/json"
	"log"
)

const PROOF_SERVICE_PATH = BASE_PATH + "present-proof/presentations"

func (c *Client) CreateProofPresentationRequest(proofReq RequestPresentationInput) (*PresentationStatus, *ApiError, int, error) {
	body, err := json.Marshal(proofReq)
	if err != nil {
		log.Fatal("Error marshalling schema: ", err)
	}
	resp, presentationStatus, apiErr, err := PostRequest[PresentationStatus](c, PROOF_SERVICE_PATH, body)
	return presentationStatus, apiErr, resp.StatusCode, err
}
func (c *Client) AcceptProofPresentation(presentation_id string) (*PresentationStatus, *ApiError, int, error) {
	body := []byte(`{ "action": "presentation-accept" }`)
	resp, presentationStatus, apiErr, err := PatchRequest[PresentationStatus](c, PROOF_SERVICE_PATH, body)
	return presentationStatus, apiErr, resp.StatusCode, err
}
func (c *Client) GetProofPresentationList() (*PresentationStatusList, *ApiError, int, error) {
	resp, presentationStatusList, apiErr, err := GetRequest[PresentationStatusList](c, PROOF_SERVICE_PATH)
	return presentationStatusList, apiErr, resp.StatusCode, err
}

func (c *Client) GetProofPresentation(presentation_id string) (*PresentationStatus, *ApiError, int, error) {
	resp, presentationStatus, apiErr, err := GetRequest[PresentationStatus](c, PROOF_SERVICE_PATH+"/"+presentation_id)
	return presentationStatus, apiErr, resp.StatusCode, err
}
