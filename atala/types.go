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
