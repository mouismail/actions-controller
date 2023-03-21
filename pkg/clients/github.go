package clients

import (
	"context"
	"fmt"
	"net/http"

	"github.tools.sap/actions-rollout-app/pkg/config"

	"github.com/bradleyfalzon/ghinstallation/v2"
	v3 "github.com/google/go-github/v50/github"
	"go.uber.org/zap"
)

type Github struct {
	logger         *zap.SugaredLogger
	keyPath        string
	appID          int64
	installationID int64
	organizationID string
	atr            *ghinstallation.AppsTransport
	itr            *ghinstallation.Transport
	serverInfo     *config.ServerInfo
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

	itr := ghinstallation.NewFromAppsTransport(atr, a.installationID)

	a.atr = atr
	a.itr = itr

	a.logger.Infow("successfully initialized github app client", "organization-id", a.organizationID, "installation-id", a.installationID, "expected-events", installation.Events)

	return nil
}

func (a *Github) VCS() config.VCSType {
	return config.Github
}

func (a *Github) Organization() string {
	return a.organizationID
}

func (a *Github) GetV3Client() *v3.Client {
	return v3.NewClient(&http.Client{Transport: a.itr})
}

func (a *Github) GetV3AppClient() *v3.Client {
	return v3.NewClient(&http.Client{Transport: a.atr})
}

func (a *Github) GitToken(ctx context.Context) (string, error) {
	t, _, err := a.GetV3AppClient().Apps.CreateInstallationToken(ctx, a.installationID, &v3.InstallationTokenOptions{})
	if err != nil {
		return "", fmt.Errorf("error creating installation token %w", err)
	}
	return t.GetToken(), nil
}
