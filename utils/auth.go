package utils

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

// AuthClient is a struct that holds the GitHub App ID, the GitHub Enterprise Server Info and a JWT token.

type ServerInfo struct {
	BaseURL   string `json:"base_url"`
	UploadURL string `json:"upload_url"`
}

type AuthClient struct {
	AppID      int64
	PrivateKey *rsa.PrivateKey
	ServerInfo ServerInfo
}

// NewAuthClient returns a new AuthClient with the given App ID and Private Key.
func NewAuthClient(appID int64, privateKeyPath string, serverInfo ServerInfo) (*AuthClient, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &AuthClient{
		AppID:      appID,
		PrivateKey: privateKey,
		ServerInfo: serverInfo,
	}, nil
}

// JWTToken returns a JWT token for the AuthClient.
func (c *AuthClient) JWTToken() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iat": now,
		"exp": now + 600, // Token is valid for 10 minutes.
		"iss": c.AppID,
	}
	token.Claims = claims

	return token.SignedString(c.PrivateKey)
}

// GitHubAppTransport is an oauth2.Transport for GitHub Apps.
type GitHubAppTransport struct {
	AuthClient *AuthClient
}

// Client returns an *http.Client that includes the GitHub App JWT token in the Authorization header.
func (t *GitHubAppTransport) Client() (*http.Client, error) {
	token, err := t.AuthClient.JWTToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return tc, nil
}

// NewGitHubClient returns a new *GitHub.Client that uses the GitHub App transport.
func (c *AuthClient) NewGitHubClient() (*github.Client, error) {
	tr := &GitHubAppTransport{AuthClient: c}
	httpClient, err := tr.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}

	client, err := github.NewEnterpriseClient(
		c.ServerInfo.BaseURL,
		c.ServerInfo.UploadURL,
		httpClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	return client, nil
}
