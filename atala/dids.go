package atala

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
