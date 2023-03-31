package clients

import (
	"context"
	"fmt"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"golang.org/x/oauth2"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	v3 "github.com/google/go-github/v50/github"
	"go.uber.org/zap"
)

type Github struct {
	logger            *zap.SugaredLogger
	keyPath           string
	appID             int64
	installationID    int64
	organizationID    string
	installationToken string
	atr               *ghinstallation.AppsTransport
	itr               *ghinstallation.Transport
	serverInfo        *config.ServerInfo
}

func NewGithub(logger *zap.SugaredLogger, organizationID string, severInfo *config.ServerInfo, config *config.GithubClient) (*Github, error) {
	a := &Github{
		logger:         logger,
		keyPath:        config.PrivateKeyCertPath,
		appID:          config.AppID,
		organizationID: organizationID,
		serverInfo:     severInfo,
	}

	err := a.initClients()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Github) initClients() error {
	atr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, a.appID, a.keyPath)
	atr.BaseURL = a.serverInfo.BaseURL
	if err != nil {
		return fmt.Errorf("error creating github app client %w", err)
	}

	enterpriseClient, err := v3.NewEnterpriseClient(a.serverInfo.BaseURL, a.serverInfo.UploadURL, &http.Client{Transport: atr})
	//installation, _, err := v3.NewClient(&http.Client{Transport: atr}).Apps.FindOrganizationInstallation(context.TODO(), a.organizationID)
	if err != nil {
		return fmt.Errorf("error creating new Enterprise Client %w", err)
	}

	installation, _, err := enterpriseClient.Apps.FindOrganizationInstallation(context.TODO(), a.organizationID)
	if err != nil {
		return fmt.Errorf("error finding organization installation %w", err)
	}

	a.installationID = installation.GetID()

	// Use the GitHub client to make API requests here
	ctx := context.Background()
	installationToken, _, err := enterpriseClient.Apps.CreateInstallationToken(ctx, a.installationID, nil)
	if err != nil {
		log.Fatalf("Error creating installation token: %v", err)
	}
	a.installationToken = installationToken.GetToken()

	newClient, err := v3.NewEnterpriseClient(a.serverInfo.BaseURL, a.serverInfo.UploadURL, oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: installationToken.GetToken()},
	)))

	if err != nil {
		log.Fatalf("Error creating new client: %v", err)
	}

	user, _, err := newClient.Users.Get(ctx, "mouismail")
	if err != nil {
		log.Fatalf("Error getting authenticated user: %v", err)
	}

	fmt.Printf("Authenticated user: %s\n", user.GetLogin())

	itr := ghinstallation.NewFromAppsTransport(atr, a.installationID)
	if err != nil {
		return fmt.Errorf("error creating github app client %w", err)
	}

	itr.BaseURL = a.serverInfo.BaseURL

	a.atr = atr
	a.itr = itr

	a.logger.Infow("successfully initialized github app client", "organization-id", a.organizationID, "installation-id", a.installationID, "expected-events", installation.Events)

	return nil
}

func (a *Github) Organization() string {
	return a.organizationID
}

func (a *Github) GetV3Client() *v3.Client {
	ctx := context.Background()
	newClient, err := v3.NewEnterpriseClient(a.serverInfo.BaseURL, a.serverInfo.UploadURL, oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.installationToken},
	)))
	if err != nil {
		a.logger.Errorw("error creating new Enterprise Client", "error", err)
	}

	return newClient
}

func (a *Github) GetV3AppClient() *v3.Client {
	client, err := v3.NewEnterpriseClient(a.serverInfo.BaseURL, a.serverInfo.UploadURL, &http.Client{Transport: a.atr})
	if err != nil {
		a.logger.Errorw("error creating new Enterprise Client", "error", err)
		return nil
	}
	return client
}

func (a *Github) GitToken(ctx context.Context) (string, error) {
	t, _, err := a.GetV3AppClient().Apps.CreateInstallationToken(ctx, a.installationID, &v3.InstallationTokenOptions{})
	if err != nil {
		return "", fmt.Errorf("error creating installation token %w", err)
	}
	return t.GetToken(), nil
}
