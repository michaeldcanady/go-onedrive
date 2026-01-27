package login

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func isEmptyRecord(r azidentity.AuthenticationRecord) bool {
	return r.ClientID == "" &&
		r.TenantID == "" &&
		r.HomeAccountID == "" &&
		r.Username == ""
}

func isAuthRequired(err error) bool {
	var authErr *azidentity.AuthenticationRequiredError
	return errors.As(err, &authErr)
}
