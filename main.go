package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/SundaeSwap-finance/atala-go/atala"
)

const (
	ISSUER_URL   = "http://localhost:8191"
	HOLDER_URL   = "http://localhost:8192"
	VERIFIER_URL = "http://localhost:8193"
)

// Parse args, create clients: issuer agent and optionally holder agent, execute specified flows
func main() {
	doAll, doInvite, doCredSchema, issueCreds, doCreateDID, doListDIDs, didRef := parseArgs()

	issuerClient := atala.CreateClient(ISSUER_URL)

	if doAll || doInvite {
		holderClient := atala.CreateClient(HOLDER_URL)
		OOB_Invitation_Flow(issuerClient, holderClient)
	}
	if doAll || doCreateDID {
		createDID_Flow(issuerClient)
	}
	if doAll || doCredSchema {
		cred_schema_flow(issuerClient)
	}
	if doAll || issueCreds {
		holderClient := atala.CreateClient(HOLDER_URL)
		issue_creds_flow(issuerClient, holderClient)
	}
	if doAll || doListDIDs {
		listDIDs(issuerClient)
	}
	if didRef != "" {
		getDID(issuerClient, didRef)
	}
}
func parseArgs() (bool, bool, bool, bool, bool, bool, string) {
	args := os.Args[1:]
	var doAll, doInvite, doCredSchema, issueCreds, doListDIDs bool
	var didRef string
	doCreateDID := len(args) == 0
	if !doCreateDID {
		switch args[0] {
		case "all":
			doAll = true
		case "invitation":
			doInvite = true
		case "doCredSchema":
			doCredSchema = true
		case "issueCreds":
			issueCreds = true
		case "getDIDs":
			doListDIDs = true
		case "getDID":
			if len(args) > 1 {
				didRef = args[1]
			} else {
				panic("getDID: Please pass the DID reference to retrieve as the second argument")
			}
		default:
			doCreateDID = true
		}
	}
	return doAll, doInvite, doCredSchema, issueCreds, doCreateDID, doListDIDs, didRef
}

// CERTIFICATE FLOW
func issue_creds_flow(issuerClient *atala.Client, holderClient *atala.Client) {
	// Get Credential Schema List
	credList, apiErr, statusCode, err := issuerClient.GetCredentialRecordsList()
	d(credList, apiErr, err, statusCode, "Existing Credential Records List:")
	// Get A Schema to Use
	var schema *atala.CredentialSchema
	schemaList := listSchemas(issuerClient)
	if len(schemaList.Contents) > 0 {
		schema = &schemaList.Contents[0]
	} else {
		// Get an Author
		var authorDid atala.DID
		dids := listDIDs(issuerClient)
		if len(dids.Contents) > 0 {
			authorDid = dids.Contents[0]
		} else {
			authorDid = *createDid(issuerClient)
		}
		schema = createCredentialSchema(issuerClient, authorDid)
	}

	var connection *atala.Connection
	connections := getConnections(issuerClient, "issuer")
	if len(connections.Contents) > 0 {
		connection = &connections.Contents[0]
	} else {
		OOB_Invitation_Flow(issuerClient, holderClient)
		connections = getConnections(issuerClient, "issuer")
		connection = &connections.Contents[0]
	}

	schemaId := ISSUER_URL + "/prism-agent/schema-registry/schemas/" + schema.Guid
	connectionId := connection.ConnectionId
	issuingDID := schema.Author // probably not correct
	// issuingDID := "did:prism:9f847f8bbb66c112f71d08ab39930d468ccbfe1e0e1d002be53d46c431212c26"
	credOfferReq := &atala.CredentialOfferRequest{
		Claims: `{
		  emailAddress": "alice@wonderland.com",
		  givenName": "Alice",
		  familyName": "Wonderland",
		  dateOfIssuance": "2020-11-13T20:20:39+00:00",
		  drivingLicenseID": "12345",
		  drivingClass": 3
		}`,
		CredentialFormat: "JWT",
		IssuingDID:       issuingDID,
		ConnectionId:     connectionId,
		SchemaId:         schemaId,
	}
	credRecord, apiErr, statusCode, err := issuerClient.CreateCredentialOffer(credOfferReq)
	d(credRecord, apiErr, err, statusCode, "Created Credential Offer:")
	recordId := credRecord.RecordId

	holderCredRecord, apiErr, statusCode, err := holderClient.GetCredentialRecord(recordId)
	d(holderCredRecord, apiErr, err, statusCode, "Holder Retrieved Credential Offer:")

	acceptedCredRecord, apiErr, statusCode, err := holderClient.AcceptCredentialOffer(recordId, credRecord.SubjectId)
	d(acceptedCredRecord, apiErr, err, statusCode, "Holder Accepted Credential Offer:")

	issuedCredRecord, apiErr, statusCode, err := issuerClient.IssueCredential(recordId)
	d(issuedCredRecord, apiErr, err, statusCode, "Issued Credential:")
}

// CREDENTIAL SCHEMA FLOW
func createCredentialSchema(issuerClient *atala.Client, authorDid atala.DID) *atala.CredentialSchema {
	// Read Schema JSON
	schema := read_schema_file("cred_schemas/sample-drivers-license-VC-schema.json")
	fmt.Println("Read Schema from File:", schema)
	credSchema := &atala.CredentialSchema{
		Name:        "driving-license",
		Version:     "1.0.0",
		Description: "Driving License Schema",
		Type:        "https://w3c-ccg.github.io/vc-json-schemas/schema/2.0/schema.json",
		Author:      authorDid.Did,
		Tags:        []string{"some tag 1", "another tag 2"},
		Schema:      *schema,
	}
	j, _ := json.MarshalIndent(credSchema, "", "    ")
	fmt.Println("Initialized Credential Schema:", string(j))
	newCredSchema, apiErr, statusCode, err := issuerClient.CreateCredentialSchema(credSchema)
	d(newCredSchema, apiErr, err, statusCode, "CREATE CERTIFICATE")
	return newCredSchema
}
func listSchemas(c *atala.Client) *atala.CredentialSchemaList {
	schemaList, _, _, _ := c.GetCredentialSchemasList()
	// schemaList, apiErr, statusCode, err := c.GetCredentialSchemasList()
	// d(schemaList, apiErr, err, statusCode, "Existing Credential Schema List:")
	return schemaList
}
func cred_schema_flow(issuerClient *atala.Client) {
	listSchemas(issuerClient)
	// Get an Author
	var authorDid atala.DID
	dids := listDIDs(issuerClient)
	if len(dids.Contents) > 0 {
		authorDid = dids.Contents[0]
	} else {
		authorDid = *createDid(issuerClient)
	}
	createCredentialSchema(issuerClient, authorDid)
	// Get Credential Schema List
	schemaList := listSchemas(issuerClient)
	lastCredSchema := schemaList.Contents[len(schemaList.Contents)-1]
	lastDotIdx := strings.LastIndex(lastCredSchema.Version, ".") + 1
	minorVersion, _ := strconv.Atoi(lastCredSchema.Version[lastDotIdx:len(lastCredSchema.Version)])
	lastCredSchema.Version = lastCredSchema.Version[:lastDotIdx] + fmt.Sprint(minorVersion+1)
	j, _ := json.MarshalIndent(lastCredSchema, "", "    ")
	fmt.Println("Changed Version of Credential Schema for Update:", string(j))
	updatedCredSchema, apiErr, statusCode, err := issuerClient.UpdateCredentialSchema(&lastCredSchema)
	d(updatedCredSchema, apiErr, err, statusCode, "UPDATE CERTIFICATE")
}
func read_schema_file(file_path string) *atala.Schema {
	content, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatal("Error while opening JSON schema file: ", err)
	}
	var payload atala.Schema
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error unmarshalling JSON schema file", err)
	}
	fmt.Printf("properties: %s\n", fmt.Sprintf("%T", payload.Properties))
	return &payload
}

// INVITATION FLOW
func OOB_Invitation_Flow(issuerClient *atala.Client, holderClient *atala.Client) {
	verifierClient := atala.CreateClient(VERIFIER_URL)
	verifierConnections := getConnections(verifierClient, "VERIFIER")
	var verifierConnection *atala.Connection
	if len(verifierConnections.Contents) > 0 {
		verifierConnection = &verifierConnections.Contents[0]
	} else {
		inv, apiErr, statusCode, err := verifierClient.CreateOOBInvitation("some label string")
		d(inv, apiErr, err, statusCode, "VERIFIER CREATE INVITATION")
		verifierConnection = inv
	}
	code := parseCodeFromInvitationURL(verifierConnection)
	fmt.Println("verifier connection", parseAsJson(verifierConnection), "\nverifier connection invitation code\n", code)
	fmt.Println("verifier connection STATE:", verifierConnection.State)

	// inv_acc, apiErr, err, statusCode := holderClient.AcceptOOBInvitation(code)
	// d(inv_acc, apiErr, err, statusCode, "HOLDER ACCEPT INVITATION")

	// inv, apiErr, err, statusCode := issuerClient.CreateOOBInvitation("some label string")
	// d(inv, apiErr, err, statusCode, "CREATE INVITATION FROM ISSUER")

	// rawInvitation := strings.SplitAfter(inv.Invitation.InvitationUrl, "https://my.domain.com/path?_oob=")[1]

	// inv_acc, apiErr, err, statusCode := holderClient.AcceptOOBInvitation(rawInvitation)
	// d(inv_acc, apiErr, err, statusCode, "ACCEPT INVITATION")

	getConnections(holderClient, "HOLDER")
	// getConnections(issuerClient, "ISSUER")
}
func parseCodeFromInvitationURL(inv *atala.Connection) string {
	return strings.SplitAfter(inv.Invitation.InvitationUrl, "https://my.domain.com/path?_oob=")[1]
}
func getConnections(c *atala.Client, who string) *atala.ConnectionList {
	cncs, apiErr, statusCode, err := c.GetConnections()
	d(cncs, apiErr, err, statusCode, who+" CONNECTIONS")
	return cncs
}

// CREATE DID FLOW
func createDID_Flow(issuerClient *atala.Client) {
	checkHealth(issuerClient)
	listDIDs(issuerClient)
	createdDid := createDid(issuerClient)
	didRef := createdDid.LongFormDid
	getDID(issuerClient, didRef)
}
func checkHealth(c *atala.Client) {
	health, apiErr, statusCode, err := c.SystemHealth()
	d(health, apiErr, err, statusCode, "HEALTH CHECK")
}
func listDIDs(c *atala.Client) *atala.DIDList {
	dids, apiErr, statusCode, err := c.ListDIDs()
	d(dids, apiErr, err, statusCode, "LIST DIDs")
	return dids
}
func createDid(c *atala.Client) *atala.DID {
	var doc = []byte(`{
		"documentTemplate": {
		  "publicKeys": [
			{
			  "id": "auth-1",
			  "purpose": "authentication"
			}
		  ],
		  "services": []
		}
	}`)
	did, apiErr, statusCode, err := c.CreateDID(doc)
	d(did, apiErr, err, statusCode, "CREATE DID")
	return did
}
func getDID(c *atala.Client, didRef string) {
	did, apiErr, statusCode, err := c.GetDID(didRef)
	d(did, apiErr, err, statusCode, "GET DID")
	didDoc, apiErr, statusCode, err := c.GetDIDDocument(didRef)
	d(didDoc, apiErr, err, statusCode, "GET DID DOCUMENT")
}

func parseAsJson(o interface{}) string {
	j, _ := json.MarshalIndent(o, "", "    ")
	return string(j)
}

// HELPER
func d(o any, apiErr *atala.ApiError, err error, statusCode int, tag string) {
	if err != nil {
		fmt.Println("###### ERROR \n######", tag, statusCode, err)
	}
	if apiErr != nil {
		fmt.Println("###### API ERROR \n######", tag, statusCode, parseAsJson(apiErr))
	}
	if o != nil {
		fmt.Println("##", tag, ":", statusCode, "\n", parseAsJson(o))
	}
}
