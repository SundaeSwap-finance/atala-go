package webhook

func NewCredentialConfig(schemaId string, credentialId string) CredentialConfig {
	return CredentialConfig{
		SchemaId:     schemaId,
		CredentialId: credentialId,
	}
}

func NewCredential(credentialId string, fields []WebhookCredentialField) WebhookCredentialValuesDTO {
	c := WebhookCredentialValuesDTO{
		CredentialId: credentialId,
		Fields:       fields,
	}
	return c
}
