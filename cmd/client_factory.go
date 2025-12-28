package cmd

import (
	"context"
	"fmt"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

func ClientFactory(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	credential, err := credentialService.LoadCredential(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %w", err)
	}

	return msgraphsdkgo.NewGraphServiceClientWithCredentials(credential, []string{"Files.ReadWrite"})
}
