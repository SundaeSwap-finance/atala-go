package atala

// Retrieve DID List
func (c *Client) ListDIDs() (*DIDList, *ApiError, int, error) {
	resp, didList, apiErr, err := GetRequest[DIDList](c, BASE_PATH+"did-registrar/dids")
	if resp == nil {
		return didList, apiErr, -1, err
	}
	return didList, apiErr, resp.StatusCode, err
}

// Retrieve DID
func (c *Client) GetDID(didRef string) (*DID, *ApiError, int, error) {
	resp, did, apiErr, err := GetRequest[DID](c, BASE_PATH+"did-registrar/dids/"+didRef)
	if resp == nil {
		return did, apiErr, -1, err
	}
	return did, apiErr, resp.StatusCode, err
}

// Retrieve resolved document for a DID
func (c *Client) GetDIDDocument(didRef string) (*DIDDocResponse, *ApiError, int, error) {
	// *NOTE: For get_DID_document endpoint: API returns status 400 and a DIDDocResponse object, instead of an APIError object. This is currently inconsistent with other endpoints.
	successCodes := []int{400}
	resp, did, apiErr, err := GetRequest[DIDDocResponse](c, BASE_PATH+"dids/"+didRef, successCodes...)
	if resp == nil {
		return did, apiErr, -1, err
	}
	return did, apiErr, resp.StatusCode, err
}

// Create a new DID with a document
func (c *Client) CreateDID(doc []byte) (*DID, *ApiError, int, error) {
	resp, createdDid, apiErr, err := PostRequest[DID](c, BASE_PATH+"did-registrar/dids", doc)
	if resp == nil {
		return createdDid, apiErr, -1, err
	}
	return createdDid, apiErr, resp.StatusCode, err
}

// DID Document Template Builder
// type PublicKey struct {
// 	id      string
// 	purpose string
// }
// type DocTemplate struct {
// 	publicKeys []PublicKey
// 	services   []string
// }
// type NewDoc struct {
// 	documentTemplate DocTemplate
// }

// func createDoc(id string, purpose string) NewDoc {
// 	var d = NewDoc{
// 		documentTemplate: DocTemplate{
// 			services: []string{},
// 			publicKeys: []PublicKey{
// 				{
// 					id:      id,
// 					purpose: purpose,
// 				},
// 			},
// 		},
// 	}
// 	return d
// }
