// identity-plugin-google provides authentication for Google accounts.
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
	"golang.org/x/oauth2/google"

	"github.com/michaeldcanady/go-onedrive/internal/core/plugins"
	identity_proto "github.com/michaeldcanady/go-onedrive/internal/features/plugins/proto/identity"
)

func mustState() string {
	uuid, err := uuid.NewV7()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate UUID for state: %v\n", err)
		return "fallback-state-" + time.Now().String()
	}

	return uuid.String()
}

type GoogleIdentityPlugin struct {
	identity_proto.UnimplementedIdentityPluginServer
}

func (p *GoogleIdentityPlugin) Login(stream identity_proto.IdentityPlugin_LoginServer) error {
	m, err := stream.Recv()
	if err != nil || m.GetConfig() == nil {
		return fmt.Errorf("expected config")
	}

	opts := m.GetConfig().Options
	clientID := opts["client_id"]
	if clientID == "" {
		return fmt.Errorf("client_id is required")
	}

	clientSecret := opts["client_secret"]
	if clientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}

	method := opts["method"]
	if method == "" {
		method = "interactive"
	}

	scopes := opts["scopes"]
	if scopes == "" {
		scopes = "https://www.googleapis.com/auth/drive https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"
	}
	scopesList := strings.Split(scopes, " ")

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopesList,
		Endpoint:     google.Endpoint,
	}

	var t *oauth2.Token
	switch method {
	case "device":
		t, err = p.loginDevice(stream, config)
	default:
		t, err = p.loginInteractive(stream, config, opts["redirect_uri"])
	}

	if err != nil {
		return err
	}

	identity, err := p.fetchIdentity(stream.Context(), config, t)
	if err != nil {
		return err
	}

	return stream.Send(&identity_proto.LoginResponse{
		Payload: &identity_proto.LoginResponse_Result{
			Result: &identity_proto.LoginResult{
				Token:    &identity_proto.AccessToken{AccessToken: t.AccessToken, RefreshToken: t.RefreshToken, ExpiresAt: t.Expiry.Unix(), Scopes: config.Scopes},
				Identity: identity,
			},
		},
	})
}

func (p *GoogleIdentityPlugin) loginInteractive(stream identity_proto.IdentityPlugin_LoginServer, config *oauth2.Config, redirect string) (*oauth2.Token, error) {
	useLocal := false
	listenAddr := "127.0.0.1:0"

	if redirect == "" || redirect == "urn:ietf:wg:oauth:2.0:oob" {
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

func (p *GoogleIdentityPlugin) loginDevice(stream identity_proto.IdentityPlugin_LoginServer, config *oauth2.Config) (*oauth2.Token, error) {
	da, err := config.DeviceAuth(stream.Context())
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("To sign in, use a web browser to open the page %s and enter the code %s to authenticate.", da.VerificationURI, da.UserCode)
	_ = p.interact(stream, &identity_proto.InteractionRequest{Action: &identity_proto.InteractionRequest_DisplayMessage{DisplayMessage: &identity_proto.DisplayMessageRequest{Message: msg}}})

	return config.DeviceAccessToken(stream.Context(), da)
}

func (p *GoogleIdentityPlugin) fetchIdentity(ctx context.Context, config *oauth2.Config, t *oauth2.Token) (*identity_proto.Identity, error) {
	resp, err := config.Client(ctx, t).Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var u struct{ Sub, Name, Email string }
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &identity_proto.Identity{Id: u.Sub, DisplayName: u.Name, Email: u.Email, Provider: "google"}, nil
}

func (p *GoogleIdentityPlugin) Refresh(ctx context.Context, req *identity_proto.RefreshRequest) (*identity_proto.RefreshResponse, error) {
	config := &oauth2.Config{ClientID: req.Options["client_id"], ClientSecret: req.Options["client_secret"], Endpoint: google.Endpoint}
	t, err := config.TokenSource(ctx, &oauth2.Token{RefreshToken: req.RefreshToken}).Token()
	if err != nil {
		return nil, err
	}
	return &identity_proto.RefreshResponse{Token: &identity_proto.AccessToken{AccessToken: t.AccessToken, RefreshToken: t.RefreshToken, ExpiresAt: t.Expiry.Unix()}}, nil
}

func (p *GoogleIdentityPlugin) interact(stream identity_proto.IdentityPlugin_LoginServer, req *identity_proto.InteractionRequest) error {
	if err := stream.Send(&identity_proto.LoginResponse{Payload: &identity_proto.LoginResponse_InteractionRequest{InteractionRequest: req}}); err != nil {
		return err
	}
	_, err := stream.Recv()
	return err
}

func (p *GoogleIdentityPlugin) GetMetadata(ctx context.Context, req *identity_proto.MetadataRequest) (*identity_proto.MetadataResponse, error) {
	return &identity_proto.MetadataResponse{
		Name:               "google",
		Type:               "identity",
		SupportedProviders: []string{"google"},
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
		Plugins:         map[string]plugin.Plugin{"identity": &plugins.IdentityGRPCPlugin{Impl: &GoogleIdentityPlugin{}}},
		GRPCServer:      plugins.CustomGRPCServer,
	})
}
