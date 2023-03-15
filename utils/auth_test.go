package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewAuthClient(t *testing.T) {
	testAppID := int64(123456)
	serverInfo := ServerInfo{
		BaseURL:   "https://api.github.com/",
		UploadURL: "https://uploads.github.com/api/uploads/",
	}
	tmpFile := createPrivateKeyForTesting()
	defer os.RemoveAll(filepath.Dir(tmpFile.Name()))
	defer tmpFile.Close()

	// Test with a valid private key file
	authClient, err := NewAuthClient(testAppID, tmpFile.Name(), serverInfo)
	if err != nil {
		t.Fatalf("failed to create AuthClient: %v", err)
	}

	if authClient.AppID != testAppID {
		t.Errorf("unexpected AppID: got %d, want %d", authClient.AppID, testAppID)
	}

	if authClient.PrivateKey == nil {
		t.Errorf("PrivateKey is nil")
	}

	// Test with an invalid private key file
	invalidFile, err := os.Create(filepath.Join("../test", "invalid-key.pem"))
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer invalidFile.Close()

	_, err = NewAuthClient(testAppID, invalidFile.Name(), serverInfo)
	if err == nil {
		t.Errorf("expected error, but got none")
	}
	if !strings.Contains(err.Error(), "failed to parse private key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAuthClientJWTToken(t *testing.T) {
	testAppID := int64(123456)
	serverInfo := ServerInfo{
		BaseURL:   "https://api.github.com/",
		UploadURL: "https://uploads.github.com/api/uploads/",
	}

	tmpFile := createPrivateKeyForTesting()
	defer os.RemoveAll(filepath.Dir(tmpFile.Name()))
	defer tmpFile.Close()

	// Create an AuthClient using the private key file.
	client, err := NewAuthClient(testAppID, tmpFile.Name(), serverInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a JWT token using the AuthClient.
	token, err := client.JWTToken()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the token is not empty.
	if token == "" {
		t.Error("generated token is empty")
	}
}

func TestGitHubAppTransport_Client(t *testing.T) {
	testAppID := int64(123456)
	serverInfo := ServerInfo{
		BaseURL:   "https://api.github.com/",
		UploadURL: "https://uploads.github.com/api/uploads/",
	}

	// Create a new GitHubAppTransport using a temporary private key file.
	tmpFile := createPrivateKeyForTesting()
	defer os.RemoveAll(filepath.Dir(tmpFile.Name()))
	defer tmpFile.Close()

	authClient, err := NewAuthClient(testAppID, tmpFile.Name(), serverInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new GitHubAppTransport using the AuthClient.
	transport := GitHubAppTransport{AuthClient: authClient}

	// Test that the transport returns a non-nil http.Client.
	client, err := transport.Client()
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Error("transport returned a nil http.Client")
	}
}

func TestNewGitHubClient(t *testing.T) {
	testAppID := int64(123456)
	serverInfo := ServerInfo{
		BaseURL:   "https://api.github.com/",
		UploadURL: "https://uploads.github.com/api/uploads//api/uploads/",
	}

	tmpFile := createPrivateKeyForTesting()
	defer os.RemoveAll(filepath.Dir(tmpFile.Name()))
	defer tmpFile.Close()

	// Create an AuthClient using the private key file.
	client, err := NewAuthClient(testAppID, tmpFile.Name(), serverInfo)
	if err != nil {
		t.Fatal(err)
	}

	// Create a GitHub client using the AuthClient.
	ghClient, err := client.NewGitHubClient()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the client is not nil.
	if ghClient == nil {
		t.Error("generated client is nil")
	}
	// Ensure the client's base URL matches the server info.
	if ghClient.BaseURL.String() != serverInfo.BaseURL {
		t.Errorf("expected client base URL %q, got %q", serverInfo.BaseURL, ghClient.BaseURL.String())
	}

	// Ensure the client's upload URL matches the server info.
	if ghClient.UploadURL.String() != serverInfo.UploadURL {
		t.Errorf("expected client upload URL %q, got %q", serverInfo.UploadURL, ghClient.UploadURL.String())
	}
}

func createPrivateKeyForTesting() *os.File {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	// Create a temporary directory and file for the private key
	tmpDir, err := os.MkdirTemp("", "github-app-test")
	if err != nil {
		log.Fatalf("failed to create temporary directory: %v", err)
	}

	keyFilePath := filepath.Join(tmpDir, "private-key.pem")
	tmpFile, err := os.Create(keyFilePath)
	if err != nil {
		log.Fatalf("failed to create temporary file: %v", err)
	}

	if err := pem.Encode(tmpFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		log.Fatal(err)
	}
	return tmpFile
}
