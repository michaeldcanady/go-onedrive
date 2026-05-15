// identity-plugin-azure provides authentication for Azure/Microsoft accounts.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"

	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
)

func mustState() string {
	uuid, err := uuid.NewV7()
	if err != nil {
		// This is highly unlikely, but we should handle it.
		// Since this is used in URLs, we'll return a random string as fallback if uuid fails,
		// but log it to stderr since we don't have a logger here yet.
		fmt.Fprintf(os.Stderr, "failed to generate UUID for state: %v\n", err)
		return "fallback-state-" + time.Now().String()
	}

	return uuid.String()
}

type AzureIdentityPlugin struct {
	identity_proto.UnimplementedIdentityPluginServer
}

func (p *AzureIdentityPlugin) Login(stream identity_proto.IdentityPlugin_LoginServer) error {
	m, err := stream.Recv()
	if err != nil || m.GetConfig() == nil {
		return fmt.Errorf("expected config")
	}

	opts := m.GetConfig().Options
	tenant := opts["tenant_id"]
	if tenant == "" {
		tenant = "common"
	}

	method := opts["method"]
	if method == "" {
		method = "interactive"
	}

	clientID := opts["client_id"]
	if clientID == "" {
		return fmt.Errorf("client_id is required")
	}

	clientSecret := opts["client_secret"]
	if clientSecret == "" && method == "client-secret" {
		return fmt.Errorf("client_secret is required for client-secret method")
	}

	redirectURI := opts["redirect_uri"]
	if redirectURI == "" {
		return fmt.Errorf("redirect_uri is required")
	}

	scopes := []string{"https://graph.microsoft.com/.default"}
	if s := opts["scopes"]; s != "" {
		scopes = strings.Split(s, ",")
	}

	var token *oauth2.Token
	switch method {
	case "device":
		token, err = p.loginDevice(stream, tenant, clientID, scopes)
	case "client-secret":
		token, err = p.loginClientSecret(stream.Context(), tenant, clientID, clientSecret, scopes)
	default:
		token, err = p.loginInteractive(stream, tenant, clientID, clientSecret, redirectURI, scopes)
	}

	if err != nil {
		return err
	}

	identity, err := p.fetchIdentity(stream.Context(), token.AccessToken)
	if err != nil {
		return err
	}

	return stream.Send(&identity_proto.LoginResponse{
		Payload: &identity_proto.LoginResponse_Result{
			Result: &identity_proto.LoginResult{
				Token: &identity_proto.AccessToken{
					AccessToken:  token.AccessToken,
					RefreshToken: token.RefreshToken,
					ExpiresAt:    token.Expiry.Unix(),
					Scopes:       scopes,
				},
				Identity: identity,
			},
		},
	})
}

func (p *AzureIdentityPlugin) loginInteractive(stream identity_proto.IdentityPlugin_LoginServer, tenant, clientID, secret, redirect string, scopes []string) (*oauth2.Token, error) {
	config := &oauth2.Config{ClientID: clientID, ClientSecret: secret, RedirectURL: redirect, Scopes: scopes, Endpoint: microsoft.AzureADEndpoint(tenant)}

	useLocal := false
	listenAddr := "127.0.0.1:0"

	if redirect == "" {
		useLocal = true
	} else if strings.HasPrefix(redirect, "http://localhost:") || strings.HasPrefix(redirect, "http://127.0.0.1:") {
		useLocal = true
		hostPort := strings.TrimPrefix(redirect, "http://")
		if colon := strings.Index(hostPort, ":"); colon != -1 {
			listenAddr = "127.0.0.1" + hostPort[colon:]
		}
	}

	if useLocal {
		l, err := net.Listen("tcp", listenAddr)
		if err != nil {
			// Fallback to random port if specified port is busy
			l, err = net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				return nil, err
			}
		}
		defer l.Close()

		// Always update redirect URI to match actual listener
		config.RedirectURL = "http://" + l.Addr().String()

		codeChan := make(chan string, 1)
		s := &http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				code := r.URL.Query().Get("code")
				if code == "" {
					return
				}
				select {
				case codeChan <- code:
					fmt.Fprintln(w, "Login successful! You can close this window.")
				default:
				}
			}),
			ReadHeaderTimeout: 10 * time.Second,
		}
		go func() {
			// nolint:staticcheck
			if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
				// We can't return an error from a goroutine, but we can log it
				// Since we don't have a logger here yet, we'll just ignore it as the timeout will catch it
			}
		}()
		defer s.Close()

		_ = p.interact(stream, &identity_proto.InteractionRequest{Action: &identity_proto.InteractionRequest_OpenUrl{OpenUrl: &identity_proto.OpenUrlRequest{Url: config.AuthCodeURL(mustState(), oauth2.AccessTypeOffline)}}})

		select {
		case code := <-codeChan:
			return config.Exchange(stream.Context(), code)
		case <-time.After(5 * time.Minute):
			return nil, fmt.Errorf("timeout")
		}
	}

	config.RedirectURL = redirect
	_ = p.interact(stream, &identity_proto.InteractionRequest{Action: &identity_proto.InteractionRequest_DisplayMessage{DisplayMessage: &identity_proto.DisplayMessageRequest{Message: "URL: " + config.AuthCodeURL(mustState(), oauth2.AccessTypeOffline)}}})
	fmt.Fprint(os.Stderr, "Code: ")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	return config.Exchange(stream.Context(), code)
}

func (p *AzureIdentityPlugin) loginDevice(stream identity_proto.IdentityPlugin_LoginServer, tenant, clientID string, scopes []string) (*oauth2.Token, error) {
	config := &oauth2.Config{ClientID: clientID, Scopes: scopes, Endpoint: microsoft.AzureADEndpoint(tenant)}
	da, err := config.DeviceAuth(stream.Context())
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate.", da.VerificationURI, da.UserCode)
	_ = p.interact(stream, &identity_proto.InteractionRequest{Action: &identity_proto.InteractionRequest_DisplayMessage{DisplayMessage: &identity_proto.DisplayMessageRequest{Message: msg}}})

	return config.DeviceAccessToken(stream.Context(), da)
}

func (p *AzureIdentityPlugin) loginClientSecret(ctx context.Context, tenant, clientID, secret string, scopes []string) (*oauth2.Token, error) {
	c := &clientcredentials.Config{ClientID: clientID, ClientSecret: secret, TokenURL: microsoft.AzureADEndpoint(tenant).TokenURL, Scopes: scopes}
	return c.Token(ctx)
}
func (p *AzureIdentityPlugin) fetchIdentity(ctx context.Context, token string) (*identity_proto.Identity, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://graph.microsoft.com/v1.0/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var me struct {
		Id                string `json:"id"`
		DisplayName       string `json:"displayName"`
		UserPrincipalName string `json:"userPrincipalName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return nil, err
	}
	return &identity_proto.Identity{Id: me.Id, DisplayName: me.DisplayName, Email: me.UserPrincipalName, Provider: "azure"}, nil
}
func (p *AzureIdentityPlugin) Refresh(ctx context.Context, req *identity_proto.RefreshRequest) (*identity_proto.RefreshResponse, error) {
	tenant := req.Options["tenant_id"]
	if tenant == "" {
		tenant = "common"
	}
	c := &oauth2.Config{ClientID: req.Options["client_id"], ClientSecret: req.Options["client_secret"], Endpoint: microsoft.AzureADEndpoint(tenant)}
	t, err := c.TokenSource(ctx, &oauth2.Token{RefreshToken: req.RefreshToken}).Token()
	if err != nil {
		return nil, err
	}
	return &identity_proto.RefreshResponse{Token: &identity_proto.AccessToken{AccessToken: t.AccessToken, RefreshToken: t.RefreshToken, ExpiresAt: t.Expiry.Unix()}}, nil
}
func (p *AzureIdentityPlugin) interact(stream identity_proto.IdentityPlugin_LoginServer, req *identity_proto.InteractionRequest) error {
	if err := stream.Send(&identity_proto.LoginResponse{Payload: &identity_proto.LoginResponse_InteractionRequest{InteractionRequest: req}}); err != nil {
		return err
	}
	_, err := stream.Recv()
	return err
}
func (p *AzureIdentityPlugin) GetMetadata(ctx context.Context, req *identity_proto.MetadataRequest) (*identity_proto.MetadataResponse, error) {
	return &identity_proto.MetadataResponse{
		Name:               "azure",
		Type:               "identity",
		SupportedProviders: []string{"azure"},
	}, nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "plugin panicked: %v\n", r)
			os.Exit(1)
		}
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugins.HandshakeConfig,
		Plugins:         map[string]plugin.Plugin{"identity": &plugins.IdentityGRPCPlugin{Impl: &AzureIdentityPlugin{}}},
		GRPCServer:      plugins.CustomGRPCServer,
	})
}
