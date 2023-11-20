package atala

import (
	"encoding/json"
	"log"
)

const SCHEMA_SERVICE_PATH = BASE_PATH + "schema-registry/"

// Get List of Credentials
func (c *Client) GetCredentialSchemasList() (*CredentialSchemaList, *ApiError, int, error) {
	resp, credSchema, apiErr, err := GetRequest[CredentialSchemaList](c, SCHEMA_SERVICE_PATH+"schemas")
	return credSchema, apiErr, resp.StatusCode, err
}

// Get Credential
func (c *Client) GetCredentialSchema(guid string) (*CredentialSchema, *ApiError, int, error) {
	resp, credSchema, apiErr, err := GetRequest[CredentialSchema](c, SCHEMA_SERVICE_PATH+"schemas/"+guid)
	return credSchema, apiErr, resp.StatusCode, err
}

// New Credential
func (c *Client) CreateCredentialSchema(schema *CredentialSchema) (*CredentialSchema, *ApiError, int, error) {
	body, err := json.Marshal(schema)
	if err != nil {
		log.Fatal("Error marshalling schema: ", err)
	}
	resp, credSchema, apiErr, err := PostRequest[CredentialSchema](c, SCHEMA_SERVICE_PATH+"schemas", body)
	return credSchema, apiErr, resp.StatusCode, err
}

// Edit Credential - saves existing doc, client must set new version
func (c *Client) UpdateCredentialSchema(schema *CredentialSchema) (*CredentialSchema, *ApiError, int, error) {
	body, err := json.Marshal(schema)
	if err != nil {
		log.Fatal("Error marshalling schema: ", err)
	}
	resp, credSchema, apiErr, err := PutRequest[CredentialSchema](c, SCHEMA_SERVICE_PATH+schema.Author+"/"+schema.Id, body)
	return credSchema, apiErr, resp.StatusCode, err
}
