package main

import (
	"atala-go/atala"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Parse args, create clients: issuer agent and optionally holder agent, execute specified flows
func main() {
	doAll, doInvite, doCreateDID, doListDIDs, didRef := parseArgs()

	issuerClient := atala.CreateClient("http://localhost:8191")

	if doAll || doInvite {
		holderClient := atala.CreateClient("http://localhost:8192")
		OOB_Invitation_Flow(issuerClient, holderClient)
	}
	if doAll || doCreateDID {
		createDID_Flow(issuerClient)
	}
	if doAll || doListDIDs {
		listDIDs(issuerClient)
	}
	if didRef != "" {
		getDID(issuerClient, didRef)
	}
}
func parseArgs() (bool, bool, bool, bool, string) {
	args := os.Args[1:]
	var doAll, doInvite, doListDIDs bool
	var didRef string
	doCreateDID := len(args) == 0
	if !doCreateDID {
		switch args[0] {
		case "all":
			doAll = true
		case "invitation":
			doInvite = true
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
	return doAll, doInvite, doCreateDID, doListDIDs, didRef
}

// INVITATION FLOW
func OOB_Invitation_Flow(issuerClient *atala.Client, holderClient *atala.Client) {
	inv, apiErr, err, statusCode := issuerClient.CreateOOBInvitation("some label string")
	d(inv, apiErr, err, statusCode, "CREATE INVITATION")

	rawInvitation := strings.SplitAfter(inv.Invitation.InvitationUrl, "https://my.domain.com/path?_oob=")[1]

	inv_acc, apiErr, err, statusCode := holderClient.AcceptOOBInvitation(rawInvitation)
	d(inv_acc, apiErr, err, statusCode, "ACCEPT INVITATION")

	getConnections(holderClient, "HOLDER")
	getConnections(issuerClient, "ISSUER")
}
func getConnections(c *atala.Client, who string) {
	cncs, apiErr, err, statusCode := c.GetConnections()
	d(cncs, apiErr, err, statusCode, who+" CONNECTIONS")
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
	health, apiErr, err, statusCode := c.SystemHealth()
	d(health, apiErr, err, statusCode, "HEALTH CHECK")
}
func listDIDs(c *atala.Client) *atala.DIDList {
	dids, apiErr, err, statusCode := c.ListDIDs()
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
	did, apiErr, err, statusCode := c.CreateDID(doc)
	d(did, apiErr, err, statusCode, "CREATE DID")
	return did
}
func getDID(c *atala.Client, didRef string) {
	did, apiErr, err, statusCode := c.GetDID(didRef)
	d(did, apiErr, err, statusCode, "GET DID")
	didDoc, apiErr, err, statusCode := c.GetDIDDocument(didRef)
	d(didDoc, apiErr, err, statusCode, "GET DID DOCUMENT")
}

// HELPER
func d(o any, apiErr *atala.ApiError, err error, statusCode int, tag string) {
	if err != nil {
		fmt.Printf("###### ERROR \n###### %s: %d\n%+v\n", tag, statusCode, err)
	}
	if apiErr != nil {
		fmt.Printf("###### API ERROR \n###### %s: %d\n%+v\n", tag, statusCode, apiErr)
	}
	if o != nil {
		j, _ := json.MarshalIndent(o, "", "    ")
		fmt.Println("##", tag, ":", statusCode, "\n", string(j))
	}
}
