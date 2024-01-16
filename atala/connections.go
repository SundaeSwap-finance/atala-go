package atala

// Retrieve Connection List
func (c *Client) GetConnections() (*ConnectionList, *ApiError, int, error) {
	resp, inv, apiErr, err := GetRequest[ConnectionList](c, BASE_PATH+"connections")
	return inv, apiErr, resp.StatusCode, err
}

// Initiate an OOB Invitation, this must be performed by the inviter
func (c *Client) CreateOOBInvitation(label string) (*Connection, *ApiError, int, error) {
	body := []byte(`{ "label": "` + label + `" }`)
	resp, inv, apiErr, err := PostRequest[Connection](c, BASE_PATH+"connections", body)
	return inv, apiErr, resp.StatusCode, err
}

// Accept an OOB Invitation, this must be performed by the invitee
func (c *Client) AcceptOOBInvitation(rawInvitation string) (*Connection, *ApiError, int, error) {
	body := []byte(`{ "invitation": "` + rawInvitation + `" }`)
	resp, inv, apiErr, err := PostRequest[Connection](c, BASE_PATH+"connection-invitations", body)
	return inv, apiErr, resp.StatusCode, err
}
