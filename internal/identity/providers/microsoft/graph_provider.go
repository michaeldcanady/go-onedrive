package microsoft

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/middleware"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	authentication "github.com/microsoft/kiota-authentication-azure-go"
	nethttp "github.com/microsoft/kiota-http-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

// GraphProvider facilitates the provisioning of an authenticated Microsoft Graph client.
type GraphProvider struct {
	// cred is the token credential used for authorizing requests.
	cred azcore.TokenCredential
	// log is the logger instance used for internal events.
	log logger.Logger
	// client is the cached Graph client instance.
	client *msgraphsdkgo.GraphServiceClient
}

// NewGraphProvider creates a new instance of GraphProvider with the provided credential and logger.
func NewGraphProvider(cred azcore.TokenCredential, log logger.Logger) *GraphProvider {
	return &GraphProvider{
		cred: cred,
		log:  log,
	}
}

// Name returns the platform identifier "microsoft".
func (p *GraphProvider) Name() string {
	return "microsoft"
}

// Adapter returns the Kiota request adapter for the Graph client.
func (p *GraphProvider) Adapter(ctx context.Context) (abstractions.RequestAdapter, error) {
	client, err := p.Client(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetAdapter(), nil
}

// Client returns an authenticated Graph client, initializing it if necessary.
func (p *GraphProvider) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	if p.client != nil {
		return p.client, nil
	}

	if p.cred == nil {
		return nil, ErrNotAuthenticated
	}

	// 1. Create the authentication provider
	authProvider, err := authentication.NewAzureIdentityAuthenticationProviderWithScopes(p.cred, []string{
		"Files.ReadWrite.All",
		"User.Read",
		"offline_access",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication provider: %w", err)
	}

	// 2. Get default middlewares and append our custom logging middleware
	handlers := nethttp.GetDefaultMiddlewares()
	handlers = append(handlers, middleware.NewKiotaLoggingMiddleware(p.log))

	// 3. Create the HTTP client with middlewares
	httpClient := nethttp.GetDefaultClient(handlers...)

	// 4. Create the request adapter
	adapter, err := nethttp.NewNetHttpRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		authProvider,
		nil,
		nil,
		httpClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request adapter: %w", err)
	}

	// 5. Create the Graph client
	client := msgraphsdkgo.NewGraphServiceClient(adapter)

	p.client = client
	return client, nil
}
