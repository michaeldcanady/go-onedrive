package auth

import (
	"time"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

type TokenResult struct {
	AccessToken string
	ExpiresOn   time.Time
	Account     public.Account
}
