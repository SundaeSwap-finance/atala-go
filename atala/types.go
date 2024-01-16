package atala

type Health struct {
	Version string
}
type ApiError struct {
	Status   int
	Type     string
	Title    string
	Detail   string
	Instance string
}

// PROOFS
type PresentationStatus struct {
	PresentationId string `json:"presentationId"`
	Thid           string `json:"thid"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	Proofs         string `json:"proofs"`
	Data           string `json:"data"`
	ConnectionId   string `json:"connectionId"`
	MetaRetries    string `json:"metaRetries"`
}
type PresentationStatusList struct {
	Contents []PresentationStatus
	Kind     string
	Self     string
	PageOf   string
	Next     string
	Prev     string
}
type RequestPresentationInput struct {
	ConnectionId     string `json:"connectionId"`
	Options          string `json:"options"`
	Proofs           string `json:"proofs"`
	CredentialFormat string `json:"credentialFormat"`
}

// CREDENTIAL OFFERS & RECORDS
type CredentialRecord struct {
	RecordId          string
	Thid              string
	CredentialFormat  string
	SubjectId         string
	ValidityPeriod    int
	Claims            string
	AutomaticIssuance bool
	CreatedAt         string
	UpdatedAt         string
	Role              string
	ProtocolState     string
	Credential        string
	IssuingDID        string
	MetaRetries       int
}
type CredentialRecordList struct {
	Contents []CredentialRecord
	Kind     string
	Self     string
	PageOf   string
	Next     string
	Prev     string
}
type CredentialOfferRequest struct {
	ValidityPeriod         int    `json:"validityPeriod"`
	SchemaId               string `json:"schemaId"`
	CredentialDefinitionId string `json:"credentialDefinitionId"`
	CredentialFormat       string `json:"credentialFormat"`
	Claims                 string `json:"claims"`
	AutomaticIssuance      bool   `json:"automaticIssuance"`
	IssuingDID             string `json:"issuingDID"`
	ConnectionId           string `json:"connectionId"`
}

// SCHEMA
type Schema struct {
	Id                   string                 `json:"$id"`
	Schema               string                 `json:"$schema"`
	Description          string                 `json:"description"`
	Type                 string                 `json:"type"`
	Properties           map[string]interface{} `json:"properties"`
	Required             []string               `json:"required"`
	AdditionalProperties bool                   `json:"additionalProperties"`
}
type Proof struct {
	Type               string
	Created            string
	VerificationMethod string
	ProofPurpose       string
	ProofValue         string
	Jws                string
	Domain             string
}
type CredentialSchemaList struct {
	Contents []CredentialSchema
	Kind     string
	Self     string
	PageOf   string
	Next     string
	Prev     string
}
type CredentialSchema struct {
	Guid        string   `json:"guid"`
	Id          string   `json:"id"`
	LongId      string   `json:"longId"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Schema      Schema   `json:"schema"`
	Author      string   `json:"author"`
	Authored    string
	Proof       Proof
	Kind        string
	Self        string
}

// DID
type DID struct {
	Did         string
	LongFormDid string
	Status      string
}
type DIDList struct {
	Contents []DID
	Kind     string
	PageOf   string
	Self     string
}

// DID Document
type DIDDocResponse struct {
	Context               string `json:"@context"`
	DidDocumentMetadata   DIDDocMeta
	DidResolutionMetadata DIDResMeta
	DidDocument           DIDDoc
}
type DIDDocMeta struct {
	Deactivated bool
	VersionId   string
}
type DIDResMeta struct {
	ContentType  string
	Profile      string
	Error        string
	ErrorMessage string
}
type DIDDoc struct {
	Context              []string `json:"@context"`
	Id                   string
	Controller           string
	Authentication       []string
	AssertionMethod      []string
	KeyAgreement         []string
	CapabilityInvocation []string
	CapabilityDelegation []string
	Service              []string
	VerificationMethod   []VerificationMethod
}
type VerificationMethod struct {
	Id           string
	Type         string
	Controller   string
	PublicKeyJwk PublicKeyJwk
}
type PublicKeyJwk struct {
	Crv string
	X   string
	Y   string
	Kty string
}

// Connections & Invitations
type Connection struct {
	ConnectionId string
	CreatedAt    string
	Kind         string
	MyDid        string
	Label        string
	Self         string
	State        string
	TheirDid     string
	UpdatedAt    string
	Invitation   Invitation
}
type Invitation struct {
	From          string
	Id            string
	InvitationUrl string
	Type          string
}
type ConnectionList struct {
	Contents []Connection
	Kind     string
	Self     string
}
