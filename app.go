package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	v3 "github.com/google/go-github/v50/github"

	"github.tools.sap/actions-rollout-app/pkg/routes"
)

var (
	webhookSecret = "development"
	appID         = int64(0) // your app id goes here
	orgID         = "your organization id"
	certPath      = "your app private cert path"

	installationID int64
	itr            *ghinstallation.Transport
)

func main() {

	appID = 61
	orgID = "mouismail-ghes"
	certPath = "test/actions-control.2023-03-17.private-key.pem"
	installationID = 199

	atr, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, appID, certPath)
	if err != nil {
		log.Fatal("error creating GitHub app client")
	}

	//installationClient, err := v3.NewEnterpriseClient(serverInfo.BaseURL, serverInfo.BaseURL, &http.Client{Transport: atr})
	//installation, _, err := installationClient.Apps.FindOrganizationInstallation(context.Background(), orgID)
	//if err != nil {
	//	log.Fatalf("error finding organization installation: %v", err)
	//}
	//
	//installationID = installation.GetID()
	itr = ghinstallation.NewFromAppsTransport(atr, installationID)

	log.Printf("successfully initialized GitHub app client, installation-id:%s expected-events:%v\n", installationID)

	http.HandleFunc("/webhook", Handle)
	http.HandleFunc("/info", routes.VersionHandler)
	http.HandleFunc("/health", routes.HealthHandler)
	err = http.ListenAndServe("0.0.0.0:3000", nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func Handle(response http.ResponseWriter, request *http.Request) {
	hook, err := ghwebhooks.New(ghwebhooks.Options.Secret(webhookSecret))
	if err != nil {
		return
	}

	payload, err := hook.Parse(request, []ghwebhooks.Event{ghwebhooks.WorkflowDispatchEvent}...)

	if err != nil {
		if err == ghwebhooks.ErrEventNotFound {
			log.Printf("received unregistered GitHub event: %v\n", err)
			response.WriteHeader(http.StatusOK)
		} else {
			log.Printf("received malformed GitHub event: %v\n", err)
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	switch payload := payload.(type) {
	case ghwebhooks.WorkflowDispatchPayload:
		log.Println("received release event")
		go processWorkflowEvent(&payload)
	case ghwebhooks.WorkflowJobPayload:
		log.Println("received workflow job event")
	case ghwebhooks.WorkflowRunPayload:
		log.Println("received workflow run event")
	default:
		log.Println("missing handler")
	}

	response.WriteHeader(http.StatusOK)
}

func GetV3Client() *v3.Client {
	client, err := v3.NewEnterpriseClient("https://github.tools.sap", "https://github.tools.sap/api/v3", &http.Client{Transport: itr})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func processReleaseEvent(p *ghwebhooks.ReleasePayload) {
	pr, _, err := GetV3Client().PullRequests.Create(context.TODO(), orgID, "target-repository", &v3.NewPullRequest{
		Title:               v3.String("Hello pull request!"),
		Head:                v3.String("develop"),
		Base:                v3.String("master"),
		Body:                v3.String("This is an automatically created PR."),
		MaintainerCanModify: v3.Bool(true),
	})
	if err != nil {
		if !strings.Contains(err.Error(), "A pull request already exists") {
			log.Printf("error creating pull request: %v\n", err)
		}
	} else {
		log.Printf("created pull request: %s", pr.GetURL())
	}
}

func processWorkflowEvent(p *ghwebhooks.WorkflowDispatchPayload) {
	ghIssue, resp, err := GetV3Client().Issues.Create(context.Background(), p.Organization.Login, p.Repository.Name, &v3.IssueRequest{
		Title:    v3.String("Hello issue!"),
		Body:     v3.String("This is an automatically created issue."),
		Assignee: v3.String("mouismail"),
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Fatalf("Creating Issue response StatusCode %s", resp.StatusCode)
	log.Printf("Issue ID: %d", ghIssue.GetID())
	//wr, err := GetV3Client().Actions.DisableWorkflowByID(context.Background(), p.Organization.Login, p.Repository.Name, p. /* workflow id goes here */)
	//if err != nil {
	//	log.Printf("error disabling workflow: %v\n", err)
	//} else {
	//	log.Printf("disabled workflow: %s", wr.Body.Read)
	//}
}
