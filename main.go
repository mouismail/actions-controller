package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"github.tools.sap/actions-rollout-app/utils"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/version", utils.VersionHandler)
	r.HandleFunc("/health", utils.HealthHandler)
	r.HandleFunc("/webhook", utils.WebhookHandler)
	r.HandleFunc("/ping", utils.HealthHandler)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3000"
	}

	log.Println("Starting server on port", httpPort)

	err := http.ListenAndServe(":"+httpPort, r)

	if err != nil {
		log.Fatalf("error occured during starting the app %s", err)
	}
}

// TODO
//func main() {
//	// Create an AuthClient with your GitHub App ID and Private Key.
//	authClient, err := githubauth.NewAuthClient(appID, privateKeyPath)
//	if err != nil {
//		log.Fatalf("failed to create AuthClient: %v", err)
//	}
//
//	// Create a new *github.Client authenticated with the AuthClient.
//	client, err := authClient.NewGitHubClient()
//	if err != nil {
//		log.Fatalf("failed to create GitHub client: %v", err)
//	}
//
//	// List the most recent events for a webhook.
//	events, _, err := client.Activity.ListRepositoryEvents(
//		context.Background(),
//		"owner",
//		"repo",
//		nil,
//	)
//	if err != nil {
//		log.Fatalf("failed to list events: %v", err)
//	}
//
//	for _, event := range events {
//		fmt.Printf("Event: %s\n", event.GetType())
//	}
//
//	// List the most recent actions runs.
//	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(
//		context.Background(),
//		"owner",
//		"repo",
//		&github.ListWorkflowRunsOptions{
//			ListOptions: github.ListOptions{PerPage: 10},
//		},
//	)
//	if err != nil {
//		log.Fatalf("failed to list actions runs: %v", err)
//	}
//
//	for _, run := range runs.WorkflowRuns {
//		fmt.Printf("Action run: %s\n", *run.HTMLURL)
//	}
//
//	// Wait for 1 minute before exiting.
//	time.Sleep(1 * time.Minute)
//}
