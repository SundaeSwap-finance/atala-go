package atala

import (
	"log"
	"net/http"
	"net/url"
)

const BASE_PATH = "prism-agent/"

// Create a new Client, pass the full url for an atala PRISM agent instance
func CreateClient(baseUrl string) *Client {
	httpClient := &http.Client{}
	c := NewClient(httpClient)
	var urlErr error
	c.BaseURL, urlErr = url.Parse(baseUrl)
	if urlErr != nil {
		log.Panicf("ERROR parsing BaseURL: %v", urlErr)
	}
	return c
}

// API Health Check
func (c *Client) SystemHealth() (*Health, *ApiError, error, int) {
	resp, health, apiErr, err := GetRequest[Health](c, BASE_PATH+"_system/health")
	return health, apiErr, err, resp.StatusCode
}

//// Connections & Invitations

// Retrieve Connection List
func (c *Client) GetConnections() (*ConnectionList, *ApiError, error, int) {
	resp, inv, apiErr, err := GetRequest[ConnectionList](c, BASE_PATH+"connections")
	return inv, apiErr, err, resp.StatusCode
}

// Initiate an OOB Invitation, this must be performed by the inviter
func (c *Client) CreateOOBInvitation(label string) (*Connection, *ApiError, error, int) {
	body := []byte(`{ "label": "` + label + `" }`)
	resp, inv, apiErr, err := PostRequest[Connection](c, BASE_PATH+"connections", body)
	return inv, apiErr, err, resp.StatusCode
}

// Accept an OOB Invitation, this must be performed by the invitee
func (c *Client) AcceptOOBInvitation(rawInvitation string) (*Connection, *ApiError, error, int) {
	body := []byte(`{ "invitation": "` + rawInvitation + `" }`)
	resp, inv, apiErr, err := PostRequest[Connection](c, BASE_PATH+"connection-invitations", body)
	return inv, apiErr, err, resp.StatusCode
}

//// DIDs

// Retrieve DID List
func (c *Client) ListDIDs() (*DIDList, *ApiError, error, int) {
	resp, didList, apiErr, err := GetRequest[DIDList](c, BASE_PATH+"did-registrar/dids")
	return didList, apiErr, err, resp.StatusCode
}

// Retrieve DID
func (c *Client) GetDID(didRef string) (*DID, *ApiError, error, int) {
	resp, did, apiErr, err := GetRequest[DID](c, BASE_PATH+"did-registrar/dids/"+didRef)
	return did, apiErr, err, resp.StatusCode
}

// Retrieve resolved document for a DID
func (c *Client) GetDIDDocument(didRef string) (*DIDDocResponse, *ApiError, error, int) {
	// *NOTE: For get_DID_document endpoint: API returns status 400 and a DIDDocResponse object, instead of an APIError object. This is currently inconsistent with other endpoints.
	successCodes := []int{400}
	resp, did, apiErr, err := GetRequest[DIDDocResponse](c, BASE_PATH+"dids/"+didRef, successCodes...)
	return did, apiErr, err, resp.StatusCode
}

// Create a new DID with a document
func (c *Client) CreateDID(doc []byte) (*DID, *ApiError, error, int) {
	resp, createdDid, apiErr, err := PostRequest[DID](c, BASE_PATH+"did-registrar/dids", doc)
	return createdDid, apiErr, err, resp.StatusCode
}

// DID Document Template Builder
type PublicKey struct {
	id      string
	purpose string
}
type DocTemplate struct {
	publicKeys []PublicKey
	services   []string
}
type NewDoc struct {
	documentTemplate DocTemplate
}

func createDoc(id string, purpose string) NewDoc {
	var d = NewDoc{
		documentTemplate: DocTemplate{
			services: []string{},
			publicKeys: []PublicKey{
				{
					id:      id,
					purpose: purpose,
				},
			},
		},
	}
	return d
}
