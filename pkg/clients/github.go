package clients

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	v3 "github.com/google/go-github/v50/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/utils"
)

type Github struct {
	logger            *zap.SugaredLogger
	keyPath           string
	appID             int64
	installationID    int64
	organizationID    string
	repository        string
	installationToken string
	atr               *ghinstallation.AppsTransport
	itr               *ghinstallation.Transport
	serverInfo        *config.ServerInfo
}

func (a *Github) GetConfig() *config.GithubClient {
	return &config.GithubClient{
		PrivateKeyCertPath: a.keyPath,
		AppID:              a.appID,
	}
}

func NewGithub(logger *zap.SugaredLogger, organizationID, repository string, severInfo *config.ServerInfo, config *config.GithubClient) (*Github, error) {
	privateKey := os.Getenv(config.PrivateKeyCertPath)
	if privateKey == "" {
		privateKey = config.PrivateKeyCertPath
	}
	a := &Github{
		logger:         logger,
		keyPath:        privateKey,
		appID:          config.AppID,
		organizationID: organizationID,
		repository:     repository,
		serverInfo:     severInfo,
	}
	err := a.initCloudClients()
	//err := a.initClients()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Github) initCloudClients() error {
	ctx := context.Background()
	atr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, a.appID, a.keyPath)

	if err != nil {
		return fmt.Errorf(utils.ErrMissingClient, err)
	}
	cloudClient := v3.NewClient(&http.Client{Transport: atr})
	installation, _, err := cloudClient.Apps.FindRepositoryInstallation(ctx, a.organizationID, a.repository)
	if err != nil {
		return fmt.Errorf(utils.ErrFindingOrgInstallations, err)
	}
	a.installationID = installation.GetID()
	a.logger.Infow("found installation id", "installation-id", a.installationID)

	installationToken, _, err := cloudClient.Apps.CreateInstallationToken(ctx, a.installationID, nil)
	if err != nil {
		return fmt.Errorf(utils.ErrCreatingInstallationToken, err)
	}
	a.installationToken = installationToken.GetToken()

	itr := ghinstallation.NewFromAppsTransport(atr, a.installationID)
	a.atr = atr
	a.itr = itr

	a.logger.Infow("successfully initialized github app client", "organization-id", a.organizationID, "installation-id", a.installationID, "expected-events", installation.Events)
	return nil
}

func (a *Github) GetCloudV3Client() *v3.Client {
	newClient := v3.NewClient(&http.Client{
		Transport: &oauth2.Transport{
			Base:   http.DefaultTransport,
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: a.installationToken}),
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	})
	return newClient
}

func (a *Github) GetCloudV3AppClient() *v3.Client {
	client := v3.NewClient(&http.Client{Transport: a.atr})
	return client
}

func (a *Github) initClients() error {
	ctx := context.Background()
	atr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, a.appID, a.keyPath)

	if err != nil {
		return fmt.Errorf(utils.ErrMissingClient, err)
	}

	atr.BaseURL = a.serverInfo.BaseURL

	enterpriseClient, err := v3.NewEnterpriseClient(a.serverInfo.BaseURL, a.serverInfo.UploadURL, &http.Client{Transport: atr})
	if err != nil {
		return fmt.Errorf(utils.ErrMissingEnterpriseClient, err)
	}

	installation, _, err := enterpriseClient.Apps.FindOrganizationInstallation(ctx, a.organizationID)
	if err != nil {
		return fmt.Errorf(utils.ErrFindingOrgInstallations, err)
	}

	a.installationID = installation.GetID()
	a.logger.Infow("found installation id", "installation-id", a.installationID)

	installationToken, _, err := enterpriseClient.Apps.CreateInstallationToken(ctx, a.installationID, nil)

	if err != nil {
		return fmt.Errorf(utils.ErrCreatingInstallationToken, err)
	}
	a.installationToken = installationToken.GetToken()

	itr := ghinstallation.NewFromAppsTransport(atr, a.installationID)
	if err != nil {
		return fmt.Errorf(utils.ErrMissingClient, err)
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

func (a *Github) Repository() string {
	return a.repository
}

func (a *Github) ServerInfo() *config.ServerInfo {
	return a.serverInfo
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
		a.logger.Errorw("error creating new Enterprise App Client", "error", err)
		return nil
	}
	return client
}
