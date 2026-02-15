package account

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func AccountToMSAuthRecord(acct Account) azidentity.AuthenticationRecord {
	return azidentity.AuthenticationRecord{
		Authority:     acct.Authority,
		ClientID:      acct.ClientID,
		HomeAccountID: acct.HomeAccountID,
		TenantID:      acct.TenantID,
		Username:      acct.Username,
		Version:       acct.Version,
	}
}

func AccountFromMSAuthRecord(rec azidentity.AuthenticationRecord) Account {
	return Account{
		Authority:     rec.Authority,
		ClientID:      rec.ClientID,
		HomeAccountID: rec.HomeAccountID,
		TenantID:      rec.TenantID,
		Username:      rec.Username,
		Version:       rec.Version,
	}
}
