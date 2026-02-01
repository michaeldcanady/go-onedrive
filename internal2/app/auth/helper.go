package auth

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func isAuthRequired(err error) bool {
	var authErr *azidentity.AuthenticationRequiredError
	return errors.As(err, &authErr)
}
