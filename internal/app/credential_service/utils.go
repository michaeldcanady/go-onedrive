package credentialservice

import "github.com/Azure/azure-sdk-for-go/sdk/azidentity"

// isEmptyRecord checks whether an AuthenticationRecord is effectively empty.
func isEmptyRecord(record azidentity.AuthenticationRecord) bool {
	return record.ClientID == "" &&
		record.TenantID == "" &&
		record.HomeAccountID == "" &&
		record.Username == ""
}
