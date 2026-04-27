package identity

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	proto "github.com/michaeldcanady/go-onedrive/internal/features/identity/proto"
)

// AuthenticateWithProto simulates calling Authenticate and mapping results to proto.
func AuthenticateWithProto(
	ctx context.Context,
	authenticator Authenticator, // e.g., *MicrosoftAuthenticator
	req *proto.AuthenticateRequest,
	log logger.Logger, // Pass logger for context
) (*proto.AuthenticateResponse, error) {
	if authenticator == nil {
		return nil, fmt.Errorf("authenticator is nil")
	}
	if req == nil {
		return nil, fmt.Errorf("nil proto.AuthenticateRequest provided")
	}

	// 1. Convert proto request to internal LoginOptions
	_, err := FromProtoAuthenticateRequest(req)
	if err != nil {
		log.Error("Failed to convert proto request to LoginOptions", logger.Error(err))
		return nil, fmt.Errorf("failed to convert proto request to LoginOptions: %w", err)
	}

	// 2. Call the actual authenticator implementation
	resp, err := authenticator.Authenticate(ctx, req)
	if err != nil {
		log.Error("Authenticator.Authenticate failed", logger.Error(err))
		return nil, fmt.Errorf("authenticator.Authenticate failed: %w", err)
	}

	// 3. Convert Go results to proto response
	// resp already contains the proto response
	log.Info("Successfully authenticated via AuthenticateWithProto")
	return resp, nil
}
