package identity

import (
	"fmt"
	"time"

	proto "github.com/michaeldcanady/go-onedrive/internal/identity/proto"
)

// --- Conversions to Proto ---

// ToProtoAuthenticateRequest converts an identity.LoginOptions struct to a proto.AuthenticateRequest.
func ToProtoAuthenticateRequest(opts LoginOptions) (*proto.AuthenticateRequest, error) {
	methodStr := opts.Method.String()
	if methodStr == "unknown" && opts.Method != AuthMethodUnknown {
		return nil, fmt.Errorf("cannot convert unknown authentication method to string for proto request")
	}

	return &proto.AuthenticateRequest{
		Method:           methodStr,
		IdentityId:       opts.AccountID,
		Force:            opts.Force,
		ProviderSpecific: opts.ProviderSpecific,
	}, nil
}

// ToProtoIdentity converts an identity.Account struct to a proto.Identity message.
func ToProtoIdentity(account Account) *proto.Identity {
	return &proto.Identity{
		Id:          account.ID,
		DisplayName: account.DisplayName,
		Email:       account.Email,
		Provider:    account.Provider,
		AvatarUrl:   account.AvatarURL,
		Metadata:    account.Metadata,
	}
}

// ToProtoAccessToken converts an identity.AccessToken struct to a proto.AccessToken message.
func ToProtoAccessToken(token AccessToken) *proto.AccessToken {
	return &proto.AccessToken{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt.Unix(),
		Scopes:       token.Scopes,
	}
}

// --- Conversions from Proto ---

// FromProtoAuthenticateRequest converts a proto.AuthenticateRequest to an identity.LoginOptions struct.
func FromProtoAuthenticateRequest(req *proto.AuthenticateRequest) (LoginOptions, error) {
	if req == nil {
		return LoginOptions{}, fmt.Errorf("nil proto.AuthenticateRequest provided")
	}

	authMethod := ParseAuthMethod(req.GetMethod())
	return LoginOptions{
		AccountID:        req.GetIdentityId(),
		Force:            req.GetForce(),
		ProviderSpecific: req.GetProviderSpecific(),
		Method:           authMethod,
		Interactive:      false,
	}, nil
}

// FromProtoIdentity converts a proto.Identity message to an identity.Account struct.
func FromProtoIdentity(identity *proto.Identity) Account {
	if identity == nil {
		return Account{}
	}
	return Account{
		ID:          identity.Id,
		DisplayName: identity.DisplayName,
		Email:       identity.Email,
		Provider:    identity.Provider,
		AvatarURL:   identity.AvatarUrl,
		Metadata:    identity.Metadata,
	}
}

// FromProtoAccessToken converts a proto.AccessToken message to an identity.AccessToken struct.
func FromProtoAccessToken(token *proto.AccessToken) AccessToken {
	if token == nil {
		return AccessToken{}
	}
	return AccessToken{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    time.Unix(token.ExpiresAt, 0),
		Scopes:       token.Scopes,
	}
}
