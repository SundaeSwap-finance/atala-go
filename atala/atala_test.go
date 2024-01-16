package atala

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"testing"
	"time"
)

const LOCAL_ATALA_HOST = "http://localhost:8191"

func TestMain(m *testing.M) {
	beforeStart(m)
	code := m.Run()
	os.Exit(code)
}

func beforeStart(m *testing.M) {
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", LOCAL_ATALA_HOST, timeout)
	if err != nil {
		fmt.Printf("\nERROR - TESTING PRE-CONDITION NOT MET:\nATALA Prism Agent not found, please restart and ensure servers are available at %s\n\n", LOCAL_ATALA_HOST)
		os.Exit(1)
	}
}

func createClient() *Client {
	issuerClient := CreateClient(LOCAL_ATALA_HOST)
	return issuerClient
}

// TestCreateDID calls atala.Client.CreateDID with a document, checks for valid return values.
func TestCreate_And_GetDID(t *testing.T) {
	name := "did:prism:"
	want := regexp.MustCompile(`^` + name)

	c := createClient()

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
	createdDid, apiErr, err, statusCode := c.CreateDID(doc)
	expectedStatusCode := 201
	if statusCode != expectedStatusCode {
		t.Fatalf("CreateDID should return %d, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	if err != nil {
		t.Fatalf("CreateDID ERROR:\n%+v\n", err)
	}
	if apiErr != nil {
		t.Fatalf("CreateDID API ERROR:\n%+v\n", apiErr)
	}
	if !want.MatchString(createdDid.LongFormDid) {
		t.Fatalf("LongFormDid: Failed to match %#q\n\nCreateDID(doc).LongFormDid = %q\n", want, createdDid.LongFormDid)
	}
	did, apiErr, err, statusCode := c.GetDIDDocument(createdDid.LongFormDid)
	if statusCode != 200 {
		t.Fatalf("After CreateDID call, GetDID should return 200, instead returned:\n%+v\n", statusCode)
	}
	if did == nil || did.DidDocument.Id != createdDid.LongFormDid {
		t.Fatalf("After CreateDID and GetDID, LongFormDids didn't match %+v != %+v", createdDid.LongFormDid, did.DidDocument.Id)
	}
}

// TestCreateDID_Error_EmptyDoc calls atala.Client.CreateDID with an empty document, checks for valid ApiError return values.
func TestCreateDID_Error_EmptyDoc(t *testing.T) {
	c := createClient()

	var doc = []byte("")
	createdDid, apiErr, err, statusCode := c.CreateDID(doc)
	expectedStatusCode := 400
	if statusCode != expectedStatusCode {
		t.Fatalf("CreateDID should return %d under API error conditions, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	if createdDid != nil {
		t.Fatalf("CreateDID should return nil object under API error conditions, instead returned:\n%+v\n", createdDid)
	}
	if err != nil {
		t.Fatalf("CreateDID should only return APIError under API error conditions, instead also returned golang error:\n%+v\n", err)
	}
	if apiErr == nil {
		t.Fatalf("CreateDID should return APIError under API error conditions")
	}
}

// TestListDIDs calls atala.Client.ListDIDs with a document, checks for valid return values.
func TestListDIDs(t *testing.T) {
	name := "did:prism:"
	want := regexp.MustCompile(`^` + name)

	c := createClient()
	dids, apiErr, err, statusCode := c.ListDIDs()
	expectedStatusCode := 200
	if statusCode != expectedStatusCode {
		t.Fatalf("ListDIDs should return %d, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	// fmt.Printf("DIDs: %+v\n", dids)
	if len(dids.Contents) > 0 {
		if !want.MatchString(dids.Contents[0].LongFormDid) {
			t.Fatalf("LongFormDid: Failed to match %#q\n\nListDIDs().Contents[0].LongFormDid = %q", want, dids.Contents[0].LongFormDid)
		}
	}
	if dids == nil || err != nil || apiErr != nil {
		t.Fatalf("CreateDID ERROR: err: %+v\napiErr:%+v", err, apiErr)
	}
}

// TestGetDID calls atala.Client.GetDID with a DID reference retrieved from ListDIDs, checks for valid return values.
func TestGetDID(t *testing.T) {
	name := "did:prism:"
	want := regexp.MustCompile(`^` + name)

	c := createClient()
	dids, apiErr, err, _ := c.ListDIDs()
	didRef := dids.Contents[0].LongFormDid
	// fmt.Printf("DIDs: %+v\n\n%+v\n\n", dids, didRef)
	// fmt.Println("DIDRef", didRef)
	did, apiErr, err, statusCode := c.GetDID(didRef)
	expectedStatusCode := 200
	if statusCode != expectedStatusCode {
		t.Fatalf("GetDID should return %d, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	// fmt.Printf("DID: %+v\n", did)
	if !want.MatchString(did.LongFormDid) {
		t.Fatalf("DidDocument.Id: Failed to match %#q\n\nGetDID().LongFormDid = %q", want, did.LongFormDid)
	}
	if did == nil || err != nil || apiErr != nil {
		t.Fatalf("GetDID ERROR: err: %+v\napiErr:%+v", err, apiErr)
	}
}
func TestGetDID_Error_BadDidRef(t *testing.T) {
	didRef := "1"
	c := createClient()
	did, apiErr, err, statusCode := c.GetDID(didRef)
	expectedStatusCode := 400
	if statusCode != expectedStatusCode {
		t.Fatalf("With bad DID ref, GetDID should return %d, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	if did != nil {
		t.Fatalf("With bad DID ref, GetDID should return nil object, instead returned:\n%+v\n", did)
	}
	if err != nil {
		t.Fatalf("With bad DID ref, GetDID should only return APIError, instead also returned golang error:\n%+v\n", err)
	}
	if apiErr == nil {
		t.Fatalf("With bad DID ref, GetDID should returned APIError")
	}
}

// TestGetDIDDocument calls atala.Client.GetDIDDocument with a DID reference retrieved from ListDIDs, checks for valid return values.
func TestGetDIDDocument(t *testing.T) {
	name := "did:prism:"
	want := regexp.MustCompile(`^` + name)

	c := createClient()
	dids, apiErr, err, _ := c.ListDIDs()
	didRef := dids.Contents[0].LongFormDid
	// fmt.Printf("DIDs: %+v\n\n%+v\n\n", dids, didRef)
	// fmt.Println("DIDRef", didRef)
	did, apiErr, err, statusCode := c.GetDIDDocument(didRef)
	// fmt.Printf("DID: %+v\n", did)
	expectedStatusCode := 200
	if statusCode != expectedStatusCode {
		t.Fatalf("GetDIDDocument should return %d, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	if !want.MatchString(did.DidDocument.Id) {
		t.Fatalf("DidDocument.Id: Failed to match %#q\n\nGetDIDDocument().DidDocument.Id = %q", want, did.DidDocument.Id)
	}
	if did == nil || err != nil || apiErr != nil {
		t.Fatalf("GetDID ERROR: err: %+v\napiErr:%+v", err, apiErr)
	}
}

// *NOTE: For get_DID_document endpoint: if API does not find doc by ref,
// it returns status 400 and a DIDDocResponse object, instead of an APIError object.
// This is currently inconsistent with other endpoints.
func TestGetDIDDocument_Error_BadDidRef(t *testing.T) {
	didRef := "1"
	c := createClient()
	did, apiErr, err, statusCode := c.GetDIDDocument(didRef)
	expectedStatusCode := 400
	if statusCode != expectedStatusCode {
		t.Fatalf("GetDIDDocument should return %d under API error conditions, instead returned:\n%+v\n", expectedStatusCode, statusCode)
	}
	if did != nil {
		if did.DidResolutionMetadata.Error != "invalidDid" {
			t.Fatalf("With bad DID ref, GetDIDDocument should return a DidResolutionMetadata.Error=invalidDid object, instead returned:\n%+v\n", did.DidResolutionMetadata.Error)
		}
	} else {
		t.Fatal("With bad DID ref, GetDIDDocument should return an object, instead returned nil")
	}
	if err != nil {
		t.Fatalf("With bad DID ref, GetDIDDocument should not return golang error, instead returned:\n%+v\n", err)
	}
	if apiErr != nil {
		t.Fatalf("With bad DID ref, GetDIDDocument should not return APIError, instead returned:\n%+v\n", apiErr)
	}
}
